# Swarm Brain

An AI based, scriptable logic system

## Concepts

### Brain

A Brain is the top level interface with a single functionality. Given a task in the form of a string it will begin calling layers to try and complete the task given.
```go
Think(task string) <-chan Result
 ```

The Brain is loaded from a set of definition files

### Layer

These are the Brain's layers. They have a single functionality, to consider a task.
```go
Consider(Signal, task.List) Signal
```

They are given one or models. For each task given the Layer will transmit, using a Transmitter, a new signal to a Neuron for processing. The resulting Signal will be added to the original Signal.

The Layers results are used as tasks for the next layer. Transmitters allow modifing how that is accomplished. By default the result(s) of the Layer is added to the task list.

### Neuron

Neurons are given a persona, prompt and a model. They have a single functionality.
```go
Work(Signal) Signal
```
When signaled they will use their persona, prompt and the given task to call an AI model. Results and errors are returned in the signal.


## Transmitter

Transmitters move signals between layers and neurons. Scripts, written in [Lua](https://www.lua.org/), can be added to modify an incoming task, an outgoing result, or to split a result into multiple new tasks.
```go
Transmit(Signal) Signal
```

### Signal

Signals carry tasks and results. Calling brain.Think() will produce one signal for it's lifetime. Many signals may be created under the primary as it works through the brain. Once completed the combined results will be returned to the brain as a Result.

## The Flow

- Ask the Brain to **Think** about a Task
- Task is converterted to a Signal
- Signal is given to the 1st Layer to **Consider**
- The Layer uses a Transmitter to **Transmit** the Signal
- The Transmitter will
	- Run any scripts for changing the task
- Give the Signal to a Neuron to do **Work**
- The Transmitter receives the Neuron result
	- Run any scripts for changing the result
	- Run any scripts that convert the result to multiple tasks
- the Layer, once all tasks are complete, will
	- Update the tasks list with the result(s)
- Signal goes to the next layer


## Defining a Brain

The swarm directory structure
```sh
|- brain.yaml
|- layers
	|- layer1.yaml
	|- layer2.yaml
|- lua
	|- script1.lua
	|- script2.lua
|- models.yaml
```

### The brain.yaml file
```yaml
---
Brain: Task Breaker
Layers:
    - tasks
    - subtasks
```

### layer files
```yaml
---
Persona: You are an expert at dividing tasks into smaller sub-tasks to help achieve a goal
Prompt: |
    Break down the task "%%TASK%%" into smaller subtasks and their expected outputs.

NeuronModels:
    - llama3
    - wizardlm2
ChangeTaskFns:
ChangeResultFns:
ResultToTasksFns:
    - split_tasks
```
Prompts can embed %%TASK%% one or more times in the prompt. If not the task is appended to the
end of the prompt.

A Layer can be allocated mutiple models. It will distribute it's task list to each model in a round robin.

### lua files
```lua
function split_tasks(str)
	local tasks = {}
	for substr in str:gmatch("%w+") do
		table.insert(tasks, substr)
  	end
	return tasks
end
```
The name of the file should match the function name.

Lua functions for changing tasks and changing results should accept a single string and return a single string.

Functions for converting a result to multiple tasks should accept a single string and return a table of strings.

### The models.yaml file
```yaml
---
- Host: OpenAI
  Model: gpt-4
  API: OpenAI
  BaseURL: https://api.openai.com/v1
  APIKey: <Your OpenAI API Key>
  Aliases:
  - gpt

- Host: Groq
  Model: llama3-70b-8192
  API: OpenAI
  BaseURL: https://api.groq.com/openai/v1
  APIKey: <Your Groq API Key>
  Aliases:
  - llama370b
  - groq

- Host: ubuntu-ai-1
  API: Ollama
  BaseURL: http://<local address>:11434
  Model: wizardlm2
  Aliases:
  - wizard

- Host: ubuntu-ai-2
  Model: llama3
  API: Ollama
  BaseURL: http://<local address>:11434
  Aliases:
  - l3
```
Models currently only support OpenAI and Ollama APIs. Model name or aliases can be used in the layer definition.
