---
layout: page
image: '/assets/img/'
description: ''
permalink: '/'
---

## What is gitsocket?

**gitsocket** is a socket (either unix socket or tcp port) server that, on
demand, will update a local repository to the latest commit of a pre-defined
repository-branch combination.

You can easily write scripts that receive webhook ([github][github-webhook] /
[bitbucket][bitbucket-webhook]) and trigger gitsocket locally. With the help of
[customized git hook] (usually bash script), this update process may further
trigger rebuild and restart of your server application.

Allow clear separation between triggering side (web server) and process side
(system user that runs git).

[githook]: https://git-scm.com/book/uz/v2/Customizing-Git-Git-Hooks

--------

## Why?

### Deploy with git push

Assuming you have your <strong>own web server</strong> hosting your web
application. Also, you're working your source code in [github][github],
[bitbucket][bitbucket] or git service alike. You see a great advantage to
deploy your source code straight from your git push action.

After some research, you've already found <strong>webhook</strong>
([github][github-webhook], [bitbucket][bitbucket-webhook]). And then you ask,
"how should I trigger my web server to git update with webhook."

It seems straightforward. All there's to do is to create a public URL that
allow github or bitbucket webhook to trigger. Then your web application will do
all the git work...wait. That doesn't sound right.

[github]: https://github.com
[bitbucket]: https://bitbucket.org
[github-webhook]: https://developer.github.com/webhooks/
[bitbucket-webhook]: https://confluence.atlassian.com/bitbucket/manage-webhooks-735643732.html

### The catch: security

To your <strong>web application</strong> update your source code with git. You'd have to:

* Allow it to access your git repository (most likely with
  <strong>ssh key</strong> that <strong>has no password</strong>)

* Allow it to access your live source code with full read / write / delete
  access.

What if the software I use has security flaw? What if an attacker hijack the
system user that my application? The attacker might:

* Stole the <strong>password-less ssh key</strong>.

* Modified your live web application and inject code to it.

### Solution: gitsocket

Naturally, you'd want to have an application receive the webhook trigger
while have another application, with all these sensitive permissions,
to do the git update on background.

**gitsocket** takes the responsibility to run all the git update process. It is
allowed to do just the relevant git command. You may apply restrictions to the
gitsocket user so it can do nothing but that. It is only responsible for the
last mile of this webhook-deploy story.

**gitsocket** expose a socket (either unix socket or tcp port) to other
application (e.g. your web application). Your application don't need to have
system permission to do the git update. It only needs to connect that socket
and trigger **gitsocket** to do it.

Now your application can trigger git update / rebuild / restart **without**
actually having the system permission to do so. Best of both world :-)

--------

## Get gitsocket

You may download the software on [github][github-gitsocket]

[github-gitsocket]: https://github.com/yookoala/gitsocket
