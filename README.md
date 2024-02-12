# AISU

AI Super U or AISU for short is a web & mobile assistant to help you set goals and achieve them.

## Overview

AISU uses a chat based language model to help you define the goals that are the most important and impactful to you and then through a combination of active encouragement, smart tracking and gamification helps you to move toward those goals with confidence.

AISU UI is written in Svelte and it's API in Go.

## Getting Started

At a minimum you'll need to install the following tools:
- Docker
- Bash

Once those are installed, you'll need to complete the configuration of the API. To do so, perform the following steps:
1. Set up your OpenAI account and get your API Key
2. Create `.env` file from the `.env.template`
3. Set the `OPENAI_API_KEY` property in the `.env` file to the value of your API key


### Running A Development Environment

To get started with the development of AISU you'll need a local environment. To this end, development containers have been defined to ensure consistency across the development environment's of anyone contributing to AISU.

_Note that executing this command mounts the current code into the devcontainer, thus any changes to the source files made from inside or outside the source files affect the current state of the code._

```bash
docker compose -f docker-compose.yaml -f docker-compose.devcontainer.yaml up -d --build  
```

This will spin up a separate dev container for the API and the UI. You can then attach to a container by simply executing 
```bash
docker container attach <container_name>
```
