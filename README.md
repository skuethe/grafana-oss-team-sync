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
    <li><a href="#setup">Setup</a></li>
      <ul>
        <li><a href="#grafana">Grafana</a></li>
        <li><a href="#entraid">EntraID</a></li>
      </ul>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
  </ul>
</details>
<br/>



<!-- ABOUT THE PROJECT -->
## About The Project

This project was created because the build-in "team sync" in Grafana is an [enterprise feature](https://grafana.com/docs/grafana/v12.0/introduction/grafana-enterprise/#team-sync) and only usable with an appropriate license key. This is where `grafana-oss-team-sync` comes into play - it is a FOSS tool to help you achive the same functionality as the "team sync" feature.  

### Current feature list

The list of features include:  

- search your enterprise solution for specific `groups` and create them as `teams` in your Grafana instance
- (optional) create `accounts` linked to each source group as `users` in your Grafana instance
- (optional) create specific `folders` in your Grafana instance and add groups to the folders `permission` list as either a `viewer`, `editor` or `admin` role

### Backlog feature list

Things which potentially will be added in the future:

- allow to reference `users` on folder permissions
- allow to reference `roles` on folder permissions
- `delete` users and groups when removed from the source / sync list

<p align="right">(<a href="#top">back to top</a>)</p>


<!-- Requirements -->
## Requirements

In it's first release, this tool only supports `Microsoft Entra ID` as a source for groups and users.  
The idea is to add new sources in the future as a "plugin" feature.

| Service  | Requirements   |
|----------|----------------|
| Grafana  | Supported versions: `11.x`, `12.x` |
| Entra ID | Auth via Azure app: `CLIENT_ID`, `TENANT_ID`, `CLIENT_SECRET` <br/>Minimum app permissions: `User.ReadBasic.All`, `GroupMember.Read.All` |

&nbsp;  
This tool is opinionated in a few ways, which result in special configuration requirements for it to work properly. See <a href="#setup">Setup</a> below.

<p align="right">(<a href="#top">back to top</a>)</p>


<!-- Setup -->
## Setup

### Grafana

soon™

### EntraID

soon™

<p align="right">(<a href="#top">back to top</a>)</p>


<!-- CONTRIBUTING -->
## Contributing

Contributions are what makes the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

See [`CONTRIBUTING`](CONTRIBUTING.md) for more information.


<p align="right">(<a href="#top">back to top</a>)</p>



<!-- LICENSE -->
## License

Distributed under the `GNU General Public License v3.0 ("GPL-3.0")`.

See [`LICENSE`](LICENSE.md) for more information.

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- CONTACT -->
## Contact

soon™

<p align="right">(<a href="#top">back to top</a>)</p>


<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
