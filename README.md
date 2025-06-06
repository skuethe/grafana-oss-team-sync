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
    <li><a href="#contact">Contact</a></li>
  </ul>
</details>
<br/>



<!-- ABOUT THE PROJECT -->
## About The Project

I created this project to get more experience with Golang and because the Grafana build-in **Team Sync** is an [enterprise feature][1] and as such only usable with an appropriate license key.  
This is where `grafana-oss-team-sync` can help you - it is a FOSS tool to help you create **teams**, **users** and even **folders** in Grafana and keep them (and their permissions) in sync with your configured source.  

Sources are setup as plugins, which can be easily extended to others in the future.  
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

In it's first release, this tool only supports `Microsoft Entra ID` as a source for groups and users.  
The idea is to add new sources in the future as a "plugin" feature.

We support Grafana versions `11.x` and `12.x`.  
It probably works with older versions as well. Just give it a try.


<!-- Configuration -->
## Configuration

The following tool specific configuration is available.  
Details on **Grafana** and **source** specific requirements can be found below.

You can configure these either in the `config.yaml` or via environment variables starting with `GOTS_`.

| Configuration                     | Config file                     | Description |
|-----------------------------------|---------------------------------|-------------|
| Log level                         | `loglevel`                      | Configure the log level<br>**Type**: `int`<br>**Env var**: `GOTS_LOGLEVEL`<br>**Allowed**: `0` (INFO), `1` (WARN), `2` (ERROR), `99` (DEBUG)<br>**Default**: `0` (INFO) |
| Source plugin                     | `source`                        | Configure the source plugin you want to use<br>**Type**: `string`<br>**Env var**: `GOTS_SOURCE`<br>**Allowed**: `entraid` |
| Auth file                         | `authFile`                      | Configure an optional file to load authentication data from<br>**Type**: `string`<br>**Env var**: `GOTS_AUTHFILE` |
| Feature: disable folder sync      | `features.disableFolders`       | Control the folder sync feature<br>**Type**: `bool`<br>**Env var**: `GOTS_FEATURE_DISABLEFOLDERS`<br>**Default**: `false` |
| Feature: disable user sync        | `features.disableUserSync`      | Control the user sync feature<br>**Type**: `bool`<br>**Env var**: `GOTS_FEATURE_DISABLEUSERSYNC`<br>**Default**: `false` |
| Feature: add local admin to teams | `features.addLocalAdminToTeams` | Control adding Grafana local admin to each team<br>**Type**: `bool`<br>**Env var**: `GOTS_FEATURE_ADDLOCALADMINTOTEAMS`<br>**Default**: `true` |
| Grafana connection                |                                 | |
|                                   | `grafana.connection.scheme`     | Configure the scheme to use<br>**Type**: `string`<br>**Env var**: `GOTS_GRAFANA_CONNECTION_SCHEME`<br>**Allowed**: `http`, `https`<br>**Default**: `http` |
|                                   | `grafana.connection.host`       | Configure the host to use<br>**Type**: `string`<br>**Env var**: `GOTS_GRAFANA_CONNECTION_HOST`<br>**Default**: `localhost:3000` |
|                                   | `grafana.connection.basePath`   | Configure the base path to use<br>**Type**: `string`<br>**Env var**: `GOTS_GRAFANA_CONNECTION_BASEPATH`<br>**Default**: `/api` |
|                                   | `grafana.connection.retry`      | Configure the connection retry, waiting 2 seconds in between each<br>**Type**: `int`<br>**Env var**: `GOTS_GRAFANA_CONNECTION_RETRY`<br>**Default**: `0` |
| Team sync                         | `teams`                         | Configure the list of teams to sync<br>**Type**: `[]string`<br>**Env var**: `GOTS_TEAMS` |
| Folder sync                       | `folders`                       | Configure the list of folders to sync<br>**Type**: `[]interface`<br>**Env var**: `GOTS_FOLDERS` |

<!-- Configuration - Grafana -->
### Grafana

| Configuration | Requirements  |
|---------------|---------------|
| Auth          | Using basic auth |
| Connection    | Modify settings in `config.yaml` path: `grafana.connection` |


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

- `Teams`: member lists of each configured team are completely overridden to avoid interference from other sources
- `Folders`: the permissions of each folder are completely overridden to avoid interference from other sources

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



<!-- CONTACT -->
## Contact

soon™

<p align="right">( <a href="#top">Back to top</a> )</p>


<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[1]: <https://grafana.com/docs/grafana/v12.0/introduction/grafana-enterprise/#team-sync> "Grafana Enterprise - Team Sync"
