# RDF4J Docker Image

<!-- spellchecker:words healthcheck -->

This image extends the tomcat variant of [eclipse/rdf4j-workbench](https://hub.docker.com/r/eclipse/rdf4j-workbench) docker image.

In particular, it adds:

- a proper healthcheck
- automatic creation of repositories on start
- proper volumes

## Environment Variables

| Variable                 | Default                              | Description                                    |
|--------------------------|--------------------------------------|------------------------------------------------|
| `RDF4J_REPOSITORY`       | `"default"`                          | Name of the repository to automatically create |
| `RDF4J_REPOSITORY_LABEL` | `"Automatically created Repository"` | Human-readable label for the repository        |

## Healthcheck

The healthcheck checks that RDF4J itself is running.
If the `RDF4J_REPOSITORY` variable is set, it checks that the given repository has been initialized.

## Volumes

- `/var/rdf4j`
- `/usr/local/tomcat/logs`
