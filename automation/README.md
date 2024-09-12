Copyright 2024 Kirk Rader

# Automation Library

- [./send_trigger_events.go](./send_trigger_events.go) exposes
  `automation.SendTriggerEvents`, a function which, as the name implies,
  launches a worker goroutine to periodically send `automation.Trigger` events
  at specific times over the course of the current day.

- [./hue/](./hue/) contains a wrapper for the V2 API exposed by Philips Hue
  Bridges, including support for receiving SSE messages asynchronously.

- [./powerview/](./powerview/) contains a wrapper for the API exposed by a
  Hunter-Douglas (PowerView) Hub for controlling motroized window shades.

- [./main/](./main/) contains a console application for manually integration
  testing all of the above.
