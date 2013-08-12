# bham - Blocky Hypertext Abstraction Markup

bham is a similar language to something like haml, but it implements
a subset of what haml does in order to keep the language neat and tidy.

----------------------

## Working Markup Examples

For a web page with the title set to 'Example Page' and an h1 tag
with the content 'Whatever' you would do the following.

```
<!DOCTYPE html>
%html
  %head
    %title Example Page
  %body
    %h1 Whatever
```

If I wanted to display a default title if there wasn't a PageTitle 
attibute in the Render Arguments when Rendering the template, and
otherwise to use the PageTitle argument, then you would do the 
following. I'm not sure that the {{ }} bit works yet...

```
<!DOCTYPE html>
%html
  %head
    = if .PageTitle
      %title
        = .PageTitle
    = else
      %title No Title Set
```

## Planned Markup Example

This following excerpt exercises most if not all of the planned 
bham features.

```
<!DOCTYPE html>
%html(ng-app)
  %head
    = $current := .Current.Page.Name
    = javascript_include_tag "jquery" "angular"
    %title Web Introduction: {{ $current }}
    = if .ExtraJsFiles
      = javascript_include_tag .ExtraJsFiles
    = stylesheet_link_tag "ui-bootstrap"
    = with $current := .Current.Variables
      = template "Layouts/CurrentJs.html" $current
  %body
    %div(class="header {{ .HeaderType }}")
    .hello Welcome to the web {{ .User.Name }}
      You are in section {{ $current }}.

    .row-fluid
      .span3
        = template "Layouts/Navigation.html" .
      .span9
        = yield
    .row-fluid
      .span10.offset1
        = range $index, $sponsor := .Sponsors
          .sponsor-mini(data-bg-image="sponsor-{{ $sponsor.Img }}")
            = link_to $sponsor.Name $sponsor.Url "class='name'"
```

## Implemented Features

* Plaintext passthrough
* %tag expansion
* If/Else Statements
* Tag Nesting
* Range statements for collection data structures
* = ... for Lines with pipelines on them

## To Be Implemented Features

* {{ }} For embedded pipeline output
* Parentheses for HTML-like attributes
* With statement for limited visibility variables
* Template Variables
* Class and ID shorthand

## Unlikely To Be Implemented Features

* Curly branch hashrocket syntax for attributes
* Multiple line prefixes for different visibility/escaping
