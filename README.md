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
        + [Full Example:](#full-example)
    * [Commands](#commands)

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

To configure your project, simply place a file named `.clonr-config.yml` into the root directory of the git repo.
(NOTE: if you are using a `.clonrrc` file -- Congratulations! You were an early adopter! `.clonrrc` will still work, but I'd suggest upgrading :D.)


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

## Commands

View documentation for the commands [here](https://github.com/OleDakotaJoe/clonr/blob/main/.resources/cmd-docs/clonr.md)





