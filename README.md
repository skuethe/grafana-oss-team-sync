<div id="top"></div>

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <!--
  <a href="https://github.com/skuethe/grafana-oss-team-sync">
    <img src="images/logo.png" alt="Logo" width="80" height="80">
  </a>
  -->
  <h1 align="center"><strong>Grafana OSS Team Sync</strong></h1>
  <p align="center">
    <a href="https://github.com/skuethe/grafana-oss-team-sync/issues">Report Bug</a>
    ·
    <a href="https://github.com/skuethe/grafana-oss-team-sync/issues">Request Feature</a>
    <br/>
    <br/>

<!-- PROJECT SHIELDS -->
<!--
*** declarations on the bottom of this document
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]
-->

  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  &nbsp;
  <ul>
    <li><a href="#about-the-project">About The Project</a></li>
    <li><a href="#requirements">Requirements</a></li>
    <li><a href="#configuration">Configuration</a></li>
      <ul>
        <li><a href="#grafana">Grafana</a></li>
        <li><a href="#source-entraid">Source: EntraID</a></li>
      </ul>
    <li><a href="#opinionated-behaviour">Opinionated Behaviour</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
  </ul>
</details>
<br/>



<!-- ABOUT THE PROJECT -->
## About The Project

I created this project to get into Go development and as such, it is probably far from being perfect. Keep an open mind to that and feel free to [contribute](#contributing) if you want to optimize or extend its functionality.  

The idea is to use `grafana-oss-team-sync` as an FOSS tool to create **teams**, **users** and even **folders** in Grafana and keep them (and their permissions) in sync with a configured source.  
This functionality _does_ exist in Grafana itself ("Team Sync"), but is a is an [enterprise feature][enterprisefeature] and as such only usable with an appropriate license key.

Sources are internally setup as plugins, which can be easily extended to others in the future.  
Currently the following sources are supported:  
- **Entra ID** (formerly "Azure Active Directory")

### Current feature list

The list of features include:  

- search your `source` for specific `groups` and create them as `teams` in your Grafana instance
- (optional) create `users` from each configured source group
- (optional) create `folders` from config input and add groups to the `permission` list as either a `viewer`, `editor` or `admin` role

### Possible future improvements

Things which potentially will be added in the future:

- allow to reference `users` on folder permissions
- allow to reference `roles` on folder permissions
- allow to assign `admin` permissions to team members
- **delete** users and groups when removed from the source / sync list

<p align="right">( <a href="#top">Back to top</a> )</p>


<!-- Requirements -->
## Requirements

In it's current state, only `Microsoft Entra ID` is available as a source for groups and users.  
The idea is to add new sources in the future as a "plugin" feature.  
Feel free to contribute your desired source.

This tool works with Grafana versions `>=11.1.0`.  


<!-- Configuration -->
## Configuration

The following tool specific configuration is available.  
Details on **Grafana** and **source** specific requirements can be found below.

You can configure these either in the `config.yaml`, via environment variables (starting with `GOTS_`) or via cli flags.  
The following hirarchy is used when merging the different config sources, overriding existing data:  
1. The `config.yaml` you specify
2. Environment variables set (also respecting a `.env` file)
3. CLI flags passed
4. (Optional) content from an `authfile`[^authfilehirarchy]

[^authfilehirarchy]: We are using [godotenv][godotenv], which will NOT override existing environment variables.  

| Configuration                     | Config via                      | Description |
|-----------------------------------|---------------------------------|-------------|
| Log level                         | **config.yaml**: `loglevel`<br>**flag**: `--loglevel` or `-l`<br>**env var**: `GOTS_LOGLEVEL` | Configure the log level<br><br>**Type**: `int`<br>**Allowed**: `0` (INFO), `1` (WARN), `2` (ERROR), `99` (DEBUG)<br>**Default**: `0` (INFO) |
| Source plugin                     | **config.yaml**: `source`<br>**flag**: `--source` or `-s`<br>**env var**: `GOTS_SOURCE`       | Configure the source plugin you want to use<br><br>**Type**: `string`<br>**Allowed**: `entraid` |
| Auth file                         | **config.yaml**: `authFile`<br>**flag**: `--authFile`<br>**env var**: `GOTS_AUTHFILE` | Configure an optional file to load authentication data from. File content needs to be in `.env` syntax (so `key=value` per line)<br><br>**Type**: `string` |
| Feature: disable folder sync      | **config.yaml**: `features.disableFolders`<br>**flag**: `--features.disableFolders`<br>**env var**: `GOTS_FEATURE_DISABLEFOLDERS` | Control the folder sync feature<br><br>**Type**: `bool`<br>**Default**: `false` |
| Feature: disable user sync        | **config.yaml**: `features.disableUserSync`<br>**flag**: `--features.disableUserSync`<br>**env var**: `GOTS_FEATURE_DISABLEUSERSYNC` | Control the user sync feature<br><br>**Type**: `bool`<br>**Default**: `false` |
| Feature: add local admin to teams | **config.yaml**: `features.addLocalAdminToTeams`<br>**flag**: `--features.addLocalAdminToTeams`<br>**env var**: `GOTS_FEATURE_ADDLOCALADMINTOTEAMS` | Control adding Grafana local admin to each team<br><br>**Type**: `bool`<br>**Default**: `true` |
| Grafana authentication            |                                 | |
| - basic auth                      | **flag**: `--username` and `--password` or `-u` and `-p`<br>**env var**: `GOTS_USERNAME` and `GOTS_PASSWORD` | Set username and password for basic authentication to Grafana<br>**Type**: `string` |
| - service account token           | **flag**: `--token` or `-t`<br>**env var**: `GOTS_TOKEN` | Set token for service account token auth to Grafana<br>**Type**: `string` |
| Grafana connection                |                                 | |
| - scheme                          | **config.yaml**: `grafana.connection.scheme`<br>**flag**: `--grafana.connection.scheme`<br>**env var**: `GOTS_GRAFANA_CONNECTION_SCHEME`     | Configure the scheme to use<br><br>**Type**: `string`<br>**Allowed**: `http`, `https`<br>**Default**: `http` |
| - host                            | **config.yaml**: `grafana.connection.host`<br>**flag**: `--grafana.connection.host` or `-h`<br>**env var**: `GOTS_GRAFANA_CONNECTION_HOST`       | Configure the host to use<br>**Type**: `string`<br>**Default**: `localhost:3000` |
| - base path                       | **config.yaml**: `grafana.connection.basepath`<br>**flag**: `--grafana.connection.basepath`<br>**env var**: `GOTS_GRAFANA_CONNECTION_BASEPATH`   | Configure the base path to use<br><br>**Type**: `string`<br>**Default**: `/api` |
| - retry                           | **config.yaml**: `grafana.connection.retry`<br>**flag**: `--grafana.connection.retry` or `-r`<br>**env var**: `GOTS_GRAFANA_CONNECTION_RETRY`      | Configure the connection retry, waiting 2 seconds in between each.<br>Only used when the return status code equals `429` or `5xx`<br><br>**Type**: `int`<br>**Default**: `0` |
| Team sync                         | **config.yaml**: `teams`<br>**flag**: `--teams` or `-t`<br>**env var**: `GOTS_TEAMS` | Configure the list of teams to sync<br><br>**Type**: `[]string` |
| Folder sync                       | **config.yaml**: `folders`<br>**flag**: `--folders` or `-f`<br>**env var**: `GOTS_FOLDERS` | Configure the list of folders to sync<br><br>**Type**: `[]interface` |

<!-- Configuration - Grafana -->
### Grafana

| Configuration | Requirements  |
|---------------|---------------|
| Version       | `>= 11.1.0` [^grafanaversion]  |
| Auth          | Using either one of the [available authentication options][availableauthoptions] `basic auth` or `service account token` [^grafanatokenauth] |

[^grafanaversion]: Minimum Grafana version is `11.1.0` as it introduced [a new bulk team membership endpoint][newbulkendpoint] we are currently using.  
[^grafanatokenauth]: Please note that `service account token` auth only works if you disable the `UserSync` feature, as creating new users in Grafana uses the Admin API, [which require the usage of basic auth][requirebasicauth].



<!-- Configuration - entraid -->
### Source: `entraid`

| Configuration   | Requirements  |
|-----------------|---------------|
| Auth            | Using Azure app via env variables: `CLIENT_ID`, `TENANT_ID`, `CLIENT_SECRET` |
| App permissions | Minimum: `User.ReadBasic.All`, `GroupMember.Read.All` |


&nbsp;  
This tool is opinionated in a few ways, which result in special configuration requirements for it to work properly. See details [below](#opinionated-behaviour).

<p align="right">( <a href="#top">Back to top</a> )</p>


<!-- Opinionated Behaviour -->
## Opinionated Behaviour

Please note the following opinionated behaviour of this tool.

- this tool should be the single point of truth for creating groups in Grafana. For that matter, we are enforcing the following:
  - `Teams`: all members of each configured team are completely overridden with matching users from the source. If you, f.e. added other additional users or canged their permission (to "admin" f.e.), these changes will be lost during the next sync operation. This also helps with keeping the groups up to date with your configured source (when removing users f.e.)
  - `Folders`: the permissions of each folder are completely overridden with the input from your config. If you don't want this to happen, you can always disable the folder sync feature via config / env variable or cli flag
- if the user sync feature is enabled, all newly created users will get a randomly generated password assigned. This password is not available afterwards, as it should not be used in the first place. Ideally you have [setup SSO authentication][setupssoauth] with the same source as your group and user sync

<p align="right">( <a href="#top">Back to top</a> )</p>


<!-- CONTRIBUTING -->
## Contributing

Contributions are what makes the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

See [`CONTRIBUTING`](CONTRIBUTING.md) for more information.


<p align="right">( <a href="#top">Back to top</a> )</p>



<!-- LICENSE -->
## License

Distributed under the `GNU General Public License v3.0 ("GPL-3.0")`.

See [`LICENSE`](LICENSE.md) for more information.

<p align="right">( <a href="#top">Back to top</a> )</p>



<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[godotenv]:             <https://github.com/joho/godotenv> "GoDotEnv"
[enterprisefeature]:    <https://grafana.com/docs/grafana/v12.0/introduction/grafana-enterprise/#team-sync> "Grafana Enterprise - Team Sync"
[availableauthoptions]: <https://grafana.com/docs/grafana/latest/developers/http_api/authentication/> "Authentication options for the HTTP API"
[newbulkendpoint]:      <https://github.com/grafana/grafana/pull/87441> "Team: Add an endpoint for bulk team membership updates"
[requirebasicauth]:     <https://grafana.com/docs/grafana/latest/developers/http_api/admin/> "Admin API"
[setupssoauth]:         <https://grafana.com/docs/grafana/next/setup-grafana/configure-security/configure-authentication/> "Configure authentication"
