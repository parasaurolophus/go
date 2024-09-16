// Copyright 2024 Kirk Rader

package hue

import (
	"parasaurolophus/utilities"
)

type (

	// Fields of interest from the /resource/scene/{id} endpoint's response.
	Scene struct {
		Name string `json:"name"`
		Id   string `json:"id"`
	}

	// Fields of interest from /resource/bridge_home, /resource/room/{id} or
	// /resource/zone/{id} endpoints' responses, plus relevant fields from
	// related scene and grouped_light resources for the given group.
	Group struct {
		Name           string           `json:"name"`
		Id             string           `json:"id"`
		Type           string           `json:"type"`
		On             bool             `json:"on"`
		GroupedLightId string           `json:"grouped_light_id"`
		Scenes         map[string]Scene `json:"scenes,omitempty"`
	}

	// Fields of interest from /resource endpoint's response, represented in a
	// way that makes them easy and efficient to use in a home automation
	// application (unlike Hue's bloated and bizarre data model),
	Model map[string]Group
)

func NewModel(resources []Item) (model Model, err error) {

	model = Model{}

	for _, resource := range resources {

		var groupType string

		if groupType, err = utilities.GetJSONPath[string](resource, "type"); err != nil {
			return
		}

		if !(groupType == "bridge_home" || groupType == "room" || groupType == "zone") {
			continue
		}

		var groupName string

		if groupName, _ = utilities.GetJSONPath[string](resource, "metadata", "name"); groupName == "" {
			groupName = "All Lights"
		}

		var groupId string

		if groupId, err = utilities.GetJSONPath[string](resource, "id"); err != nil {
			return
		}

		var (
			group Group
			ok    bool
		)

		if group, ok = model[groupName]; !ok {

			var (
				resourceType, ownerId, resourceId string
				groupedLightState                 bool
			)

			for _, r := range resources {

				if resourceType, err = utilities.GetJSONPath[string](r, "type"); err != nil {
					return
				}

				if resourceType != "grouped_light" {
					continue
				}

				if ownerId, err = utilities.GetJSONPath[string](r, "owner", "rid"); err != nil {
					return
				}

				if ownerId != groupId {
					continue
				}

				if resourceId, err = utilities.GetJSONPath[string](r, "id"); err != nil {
					return
				}

				if groupedLightState, err = utilities.GetJSONPath[bool](r, "on", "on"); err != nil {
					return
				}

				break
			}

			group = Group{
				Name:           groupName,
				Id:             groupId,
				Type:           groupType,
				GroupedLightId: resourceId,
				On:             groupedLightState,
				Scenes:         map[string]Scene{},
			}

			for _, r := range resources {

				if resourceType, err = utilities.GetJSONPath[string](r, "type"); err != nil {
					return
				}

				if resourceType != "scene" {
					continue
				}

				if ownerId, err = utilities.GetJSONPath[string](r, "group", "rid"); err != nil {
					return
				}

				if ownerId != groupId {
					continue
				}

				var sceneName, sceneId string

				if sceneName, err = utilities.GetJSONPath[string](r, "metadata", "name"); err != nil {
					return
				}

				if sceneId, err = utilities.GetJSONPath[string](r, "id"); err != nil {
					return
				}

				scene := Scene{
					Name: sceneName,
					Id:   sceneId,
				}

				group.Scenes[sceneName] = scene
			}
		}

		model[groupName] = group
	}

	return
}

func (group Group) GetSceneById(id string) (scene Scene, ok bool) {

	for _, s := range group.Scenes {

		if s.Id == id {

			scene = s
			ok = true
			return
		}
	}

	return
}

func (model Model) GetGroupById(id string) (group Group, ok bool) {

	for _, g := range model {

		if g.Id == id {
			group = g
			ok = true
			return
		}
	}

	return
}

func (model Model) GetScene(groupName string, sceneName string) (scene Scene, ok bool) {

	var group Group

	if group, ok = model[groupName]; !ok {

		return
	}

	scene, ok = group.Scenes[sceneName]
	return
}
