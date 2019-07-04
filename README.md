# git-follow-up

git-follow-up aims at keeping track of contributions made on multiple git repositories described in a yaml configuration file.
Those repositories can be hosted on any platform, and accessed through ssh, https, with or without an access token.

Each git project is synced concurrently as a bare repository locally, in `~/.git-follow-up/git/` directory.

Then the commits list can be queried and filtered with the provided flags.  

![Screenshot of git-follow-up](https://user-images.githubusercontent.com/1558361/60553906-b2cb3480-9d35-11e9-9e4b-fef28cecd9ec.png)

### Why this project?

I've been working in organizations where there are a lot of git repos (microservices, packaging, config, ...) and I wasn't fully satisfied with the git platform web ui. 
My need was simply to check what happened on all repos, so as to quickly prepare a stand-up meeting, or keep track of the overall progress.
Therefore, I chose to develop my own tool in golang, using the occasion to learn this language ;)

### Installation

If you have go already installed, you can run : 
```bash
go get github.com/ttauveron/git-follow-up
```

You can also download the binary matching your OS (https://github.com/ttauveron/git-follow-up/releases) and copy it, for example, in `/usr/local/bin` (don't forget to `chmod +x` the binary).

### Configuration

The `config.yaml`, that should be located by default in the `~/.git-follow-up/` directory, references all the repositories that need to be tracked.

Here is an example of that config file describing the repositories : 

```yaml
repositories:
  - name: go-git
    url: git@github.com:src-d/go-git.git
    authentication:
      type: ssh
      auth_file: /home/ttauveron/.ssh/id_rsa
    labels:
      - go
      - git
      
  - name: cobra
    url: https://github.com/spf13/cobra.git
    authentication:
      type: access_token
      auth_file: /home/ttauveron/.git-follow-up/gh_access_token
    labels:
      - go
      
  - name: viper
    url: https://github.com/spf13/viper
```

#### Description of the yaml fields

| Field name | Description |
|------|-------------------------------|
|name  | the name given to the project |
| url | url of the git repo (ssh, https)| 
| authentication | The types available are *ssh* and *access_token*. <br>  The *auth_file* parameter specifies the key to be used to authenticate to the git hosting platform you're using. <br> For a ssh authentication, we are pointing to a ssh private key file and for a https authentication, we are pointing to a file containing the access token provided by the git hosting platform.| 
|labels| Labels add filtering options to repositories, allowing to query a subset of the defined repositories |

### Usage

First, we need to sync the repositories defined in the config file locally.

(Note : for large repos, this step may initially take some time...)

```bash
git-follow-up update 
```

Then we can query the local repositories for commits
```bash
git-follow-up commits --from 2019-01-10 --author ttau --label go --label git
```

The commit command can take multiple flags : 

| Flags| Description| 
|---|---| 
|--from| Filters commit by date<br>Default value : "wtd" (week to date) <br><br> Possible values : <br>- today<br>- yesterday<br>- wtd<br>- mtd<br>- ytd<br>- yyyy-MM-dd|
|--author| Filters commit by author <br>This flag can be specified multiple times for targeting multiple authors|
|--display|Commit fields to be displayed (all by default)<br>This flag can be specified multiple times for displaying multiple fields<br><br>Possible values :<br>- author<br>- date<br>- hash<br>- message<br>- repo|  
|--label|Filters by project labels<br>This flag can be specified multiple times to target multiple labels|
|--update|Runs the update command before querying the repos|

For example, we can list contributors on a time range : 
```bash
git-follow-up commits --from ytd --display author | sort | uniq
```

### Bash completion

To activate bash completion for git-follow-up, run the following command :

```bash
source <(git-follow-up completion)
```

Alternatively, add this command to your ~/.bashrc file to persist the bash completion.
