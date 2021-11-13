# Clonr Project Templating CLI

- [About](#about)
    * [Installation](#installation)
        + [Homebrew](#homebrew)
        + [Go install](#go-install)
        + [npm](#npm)
    * [Quick start for developers](#quick-start-for-developers)
    * [Configuring a project.](#configuring-a-project)
        + [Basic Example](#basic-example)
        + [Example With Globals](#example-with-globals)
        + [Validation](#validation)
        + [Conditionals](#conditionals)
            - [Clonr's Runtime Data Transfer API](#clonrs-runtime-data-transfer-api)
                * [`getClonrVar()`](#getclonrvar)
                * [`getClonrBool()`](#getclonrbool)
                * [`clonrResult`](#clonrresult)
            - [Conditional File Rendering](#conditional-file-rendering)
                * [Single-line script](#single-line-script)
                * [Multi-line script](#multi-line-script)
            - [Conditional Text Block](#conditional-text-block)
            - [Best Practices](#best-practices)
            - [Supported Javascript Syntax](#supported-javascript-syntax)
        + [Full Example:](#full-example)
        + [Using Aliases](#using-aliases)
    * [Commands](#commands)(base)


# About
This project is aimed to make creating template projects very easy, so that you can set up a project one time, and not worry about configuration again.
Simply host your template project in a git repostory, configure your template variables in a .clonr-config.yml file, as well as providing a placeholder in the
template files, and run `clonr clone <repo_url>`. The rest will unfold before your eyes.

## Installation

### Homebrew
If you would like to install the project via homebrew:
`brew install oledakotajoe/clonr/clonr`
then  run
`clonr version`
to check the installation

### Go install

If you have go installed on your machine

`go install github.com/oledakotajoe/clonr`

Alternatively you may clone this repo, then run `go build && go install`
(if you already have your PATH set up for golang, then you should be good to go)

### npm

```shell
npm install -g go-clonr
go-clonr install
clonr version
```

If you are on a unix/linux machine get access denied exception when running this command, you can either choose a
different installation method, or run sudo

```shell
sudo npm install -g go-clonr
clonr version
```

If this still does not work, try a different method of installation.

## Quick start for developers
Make sure you have Go installed on your machine. [Find out how](https://golang.org/doc/install)

In your terminal:
1. Clone the project
   `git clone https://github.com/OleDakotaJoe/clonr.git`
2. cd into the projects directory
3. Run `go build` to download all dependencies
4. Run `go run main.go version` to verify the install.

Now you can use the CLI!

To play with an example project, open your terminal and run

`go run main.go clone https://github.com/OleDakotaJoe/clonr-example-template.git`

This will create a copy of the above git repo on your local machine under the directory 'clonr-app'.
The CLI will always install your project in your present working directory.
After you run this command, any template variables that are configured in your .clonr-config.yml file will be picked up by the engine,
and you will be asked to provide input via the terminal.

## Configuring a project.

To configure your project, place a file named `.clonr-config.yml` into the root directory of the git repo.
(NOTE: if you are using a `.clonrrc` file -- Congratulations! You were an early adopter! `.clonrrc` will still work, but I'd suggest upgrading :D )


Inside this yaml file you will need to provide the paths to the files which need to be processed, the name of the
placeholder variables that you have provided in those files, and the questions which need to be asked to determine those variables.

### Basic Example
The root key in the yaml file must be "templates"

Here is an example of what the yaml syntax for clonr looks like:
```yaml
templates:
  README.md:
    location: /README.md
    variables:
      clonr_variable:
        question: What do you want the value of this variable to be?
  LICENSE:
    location: /LICENSE
    variables:
      owner:
        question: Who is the owner of the project?
      date:
        question: When did your Copyright begin?
```

Variables can only contain lowercase letters, numbers, dashes, and underscores.

Syntax for placeholder variables within the template files:
```
{@clonr{your_variable_inside_these_brackets}}
```

You can include as many of these inside the files as you would like.


### Example With Globals
You may declare global variables by providing the key: `globals`. This will allow you to ask a question only one time,
and use that variable in more than one file, without being required to get input from the user again.
Under globals, provide the key named `variables` and then simply provide the key for each variable name, and include your question below that.


In order to consume a global variable that you declared, in the variables section of the `templates` key, include the key `globals`.
you do not need to provide a value for that key, clonr will search the file for any placeholders that match a global variable, and make the swap.

NOTE: If you do not need a variable to be available to more than one file, then provide that variable under templates, or you may see some performance issues.

Here is an example of a use case:
```yaml
globals:
  variables:
    project-name:
      question: What do you want the project name to be?
templates:
  some-file.txt:
    location: /some-file.txt
    variables:
      globals: # This key lets clonr know that there are global variables in /some-file.txt, and to scan for them. You do not need to provide a corresponding value for this key
      some-other-variable:
        question: What do you want the value of some-other-variable to be?
```

The syntax for the placeholder for a global variable is as follows:
```
{@clonr{globals.project-name}}
```
Note that the syntax is identical, EXCEPT prefix your variable name with `globals.`

### Validation
It is recommended that you provide a regex for validating the end-user input for your template variables. 
To add validation, simply add a `validation` key under the variable you would like to validate, and pass in a regex 
for validation. 
Be aware that the regex will be evaluated by golang, and must follow the [go regex spec](https://yourbasic.org/golang/regexp-cheat-sheet/).

NOTE: If a default variable is provided, it is not checked against validation in this step.

Validation Example: 
```yaml
templates: 
  some-file: 
    location: /some-file.txt
    variables:
      some-variable:
        question: you got this part by now.
        validation: "[\\w]" // This would correspond to a "word" type regex.
```

It is important that you escape your slashes and wrap in double quotes like seen above. 
For example the regex: `[\w]` should be provided in this syntax: `"[\\w]"`
Note the two \'s and "s.
For more information on golang regexes, see [this cheat sheet](https://yourbasic.org/golang/regexp-cheat-sheet/)
For more information on escaping regexes see [this article](https://www.threesl.com/blog/special-characters-regular-expressions-escape/)

### Conditionals
If you need to conditionally render a block of text, or an entire file, you can use clonr's built-in javascript runtime
to process your logic and conditionally render text.

*Clonr provides a simple API for resolving a script written in JS.*


#### Clonr's Runtime Data Transfer API
If you want to get the value of a variable (you can access any variable defined in the `.clonr-config.yml` file)
these methods are available:

##### `getClonrVar()`
To use `getClonrVar()` pass in a string that matches this syntax: 
```"<template>[<variable>]"```
In this example replace `<template>` with either `globals` or the corresponding template name that you've chosen.
(Note: do not use the full filepath of the template, just use the key-name of the template). 
Replace `<variable>` with the variable that belongs to that template, or the global variable name. When using global variables 


Example in context:
```
{@clonr{%
    const switchCase = getClonrVar("globals[some-var]") // This line will get the value that was provided for the global variable "some-var"
    switch (switchCase) {
    case "some-choice": // if the user had provided "some-choice" as the value for "some-var", this case would be executed
        clonrResult = "whatever"
        break;
    case "some-other-choice": // if the user had provided "some-other-choice" as the value for "some-var", this case would be executed
        clonrResult = "something-else"
        break;
    }
%}/clonr}
```

##### `getClonrBool()`
This method works exactly the same as `getClonrVar()` except it explicitly returns a boolean, and not a string.
This is based on Golang's casting system, not javascript's.

##### `clonrResult`
This is a protected variable. 
Worth mentioning, you do not have to "initialize" `clonrResult`. You can simply set its value.
Whether you are conditionally rendering a file, or conditionally rendering a block of text, you must set the result of 
your logic to `clonrResult`. 

In the case of conditionally rendering a block of text, the resulting value must be a string. 
In the case of a conditionally rendering an entire file, the resulting value must be boolean.

#### Conditional File Rendering
If you want an entire file to be conditional, define a conditional block inside your `.clonr-config.yml` file, 
under the template that you want to be conditional.

Note that when using conditional file rendering, you MUST set clonrResult to a boolean value, or you may get unwanted
results.

##### Single-line script
Example:
```yaml
templates:
  some-file: 
    location: /some-file.txt
    variables:
      some-variable:
        question: True or False?
        choices:
          - true
          - false
    condition: clonrResult = getClonrBool("some-file[some-variable]")
```
When running clonr in this scenario the user will be asked `True or False?` and given a choice. 
If the user chooses `true` the file will be rendered
If the user chooses `false` the file will not be rendered

##### Multi-line script
If you need to use a multiline script, use yaml's [string-literal-format](https://symfony.com/doc/current/components/yaml/yaml_format.html#strings)
To use the string literal format, follow your `condition` key with a pipe character `|` then a new line.
Like this: 
```yaml
condition: |
  // your multi
  // line script
  // here.
```

Example in context: 

```yaml
templates:
  some-file: 
    location: /some-file.txt
    variables:
      some-variable:
        question: Make a decision!
        choices:
          - some-enum
          - some-other-enum
    condition: | 
      if (getClonrVar("some-file[some-variable]") === "some-enum") {
        clonrResult = true
      } else {
        clonrResult = false
      }
```

#### Conditional Text Block
If you want to conditionally render a text block use the following syntax:
```
{@clonr{%
    if (getClonrBool("globals[some-var]")) {
        clonrResult = "some-value"
    } else {
        clonrResult = ""     
    }
%}/clonr}
```

In this scenario, everything between `{@clonr{%` and `%}/clonr}` will be executed in the javascript runtime.
If the value returned by `getClonrBool("globals[some-var]")` is truthy, then this block will render "some-value". If the
value was falsy this block will render nothing.

#### Best Practices
It is a best practice to use an enum (choices) wherever possible when dealing with string based conditionals. 
This will allow you to ensure that the input the user provides is acceptable for your logic.

#### Supported Javascript Syntax
The runtime is handled by [goja](https://github.com/dop251/goja), and uses ES5 syntax.
This runtime will also support some ECMAScript6 syntax, for example: arrow functions, const. You can find a full list of the 
supported ES6 syntax [here](https://github.com/dop251/goja/milestone/1?closed=1)

Also be aware, that the ECMAScript standard is not the same as browser APIs, and some functionality you are used to in 
the browser is not going to be available in this application. Example: `fetch` `WebGL`, anything `DOM` related, File System Access API. 
You can find a more comprehensive list of browser API's [here](https://developer.mozilla.org/en-US/docs/Web/API)
If you are trying to run a script, and getting a message indicating that your method is not defined, check the above list. 

### Full Example:
```yaml
globals:
  variables:
    project-name:
      question: What do you want the project name to be?
templates:
  package.json:
    location: /package.json
    variables:
      globals: # This key lets clonr know that there are global variables in /some-file.txt, and to scan for them. You do not need to provide a corresponding value for this key
      starting-version:
        question: What do you want the starting version to be?
```

And in your package.json file you might have this:
```json
{
  "name": "{@clonr{globals.project-name}}",
  "version": "{@clonr{starting-version}}"
}
```

When you run the CLI you will be asked:

```
    What do you want the package name to be?
```

Type in your response.
Let's say for example, my response was `awesome-react-app

Then you'll be asked
```
What do you want the starting version to be?
```
Type in your response. Lets say I said "1.0.0"

Go check your files:
Your file should now look like this
```json
{
  "name": "awesome-react-app",
  "version": "1.0.0"
}
```

### Using Aliases
You can set aliases for your projects by using the `clonr alias [args] [flags]` command. 
See more information about aliases [here](https://github.com/OleDakotaJoe/clonr/blob/main/.resources/cmd-docs/clonr_alias.md)

## Commands

View documentation for the commands [here](https://github.com/OleDakotaJoe/clonr/blob/main/.resources/cmd-docs/clonr.md)


