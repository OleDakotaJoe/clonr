# Clonr Project Templating CLI

# About
This project is aimed to make creating template projects very easy, so that you can set up a project one time, and not worry about configuration again. 
Simply host your template project in a git repostory, configure your template variables in a clonr.yaml file, as well as providing a placeholder in the 
template files, and run `clonr clone <repo_url>`. The rest will unfold before your eyes.

## Quick start
Make sure you have Go installed on your machine. [Find out how](https://golang.org/doc/install)

In your terminal: 
1. Clone the project 
`git clone https://github.com/OleDakotaJoe/clonr.git`
2. cd into the projects directory
3. Run `go get` to download all dependencies
4. Run `go install` to install the project.
5. Run `clonr version` to verify the install.

Now you can use the CLI!

To play with an example project, open your terminal and run 

`clonr clone https://github.com/OleDakotaJoe/clonr-example-template.git`

This will create a copy of the above git repo on your local machine under the directory 'clonr-app'.
The CLI will always install your project in your present working directory.
After you run this command, any template variables that are configured in your clonr.yaml file will be picked up by the engine, 
and you will be asked to provide input via the terminal.

## Configuring a project.

To configure your project, simply place a file named `clonr.yaml` into the root directory of the git repo. 
Inside this yaml file you will need to provide the paths to the files which need to be processed, the name of the 
placeholder variables that you have provided in those files, and the questions which need to be asked to determine those variables.

The root key in the yaml file must be "paths"

Here is an example of what a yaml file looks like:
```yaml
paths:
  README.md:
    location: /README.md
    variables:
      - clonr_variable: What do you want the value of this variable to be?
  LICENSE:
    location: /LICENSE
    variables:
      - owner: Who is the owner of the project?
      - date: When did your Copyright begin?
```

Variables can only contain lowercase letters, numbers, dashes, and underscores.

Syntax for placeholder variables within the template files:
```
{@clonr{your_variable_inside_these_brackets}}
```

You can include as many of these inside the files as you would like.

EXAMPLE: 
```yaml 
paths: 
  package.json:
    location: /package.json
    varialbes:
      - package-name: What do you want the package name to be?
```

And in your package.json file you might have this:
```json
{
  "name": "{@clonr{package-name}}",
  "version": "1.0.0"
}
```

When you run the CLI you will be asked:

```
    What do you want the package name to be?
```

Type in your response and then go check your results.
Let's say for example, my response was `awesome-react-app`
Your file should now look like this
```json
{
  "name": "awesome-react-app",
  "version": "1.0.0"
}
```


## Known Issues
Currently only https git addresses are supported
