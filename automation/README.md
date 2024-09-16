Copyright 2024 Kirk Rader

# Automation Library

- [./trigger/](./trigger/) exposes `automation.NewTimer`, a function which, as
  the name implies, launches a worker goroutine to periodically send
  `automation.Trigger` events at specific times over the course of the current
  day. The goroutine and the channels necessary interoperate with it are
  wrapped in the `trigger.Timer` struct, but note that all the actual behavior
  is encapsulated in the body of the `trigger.NewTimer` function, itself.

- [./hue/](./hue/) contains a wrapper for the V2 API exposed by Philips Hue
  Bridges, including support for receiving SSE messages asynchronously.

- [./powerview/](./powerview/) contains a wrapper for the API exposed by a
  Hunter-Douglas (PowerView) Hub for controlling motroized window shades.

## Requirements

This began as an idle experiment to see what it would be like to replace any
part of my fairly elaborate, Node-RED based [home automation
system](https://github.com/parasaurolophus/home-automation) with code written
in Go. I began by focusing on idiomatic Go implementations of the core back end
features of the existing system:

- Access to the JSON-based data models provided by some specific "smart" device
  manufacturers' API's (Philips Hue and Hunter-Douglas PowerView)

- Subscribe to asynchronous SSE messages that are part of one such API (Philips Hue)

- Send asynchronous messages at specific times of day to drive lighting and
  window covering automation

## General Principles

By "idiomatic" Go code, I mean:

- Use goroutines and channels for asynchronous operations

- Rely on communication via channels rather than access to shared memory in
  goroutines to the greatest degree practical

