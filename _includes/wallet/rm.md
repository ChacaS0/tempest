
<a href="https://asciinema.org/a/173451" target="_blank"><img src="https://asciinema.org/a/173451.png" /></a>

### [](#head-rm-1)List of flags

| Long Flag       | Short Flag  | Description                                                                                              |
|:----------------|:------------|:---------------------------------------------------------------------------------------------------------|
| \-\-help        | -h          | [NO_ARGS] See the help section for the command.                                                          |
| \-\-origin      | -o          | [NO_ARGS] To be used with ``-p`` or ``-i``, indicates to **TEMPest** to remove it from your computer.    |


### [](#head-rm-2)Explanations

| Example command                      | Description                                                                                       |
|:-------------------------------------|:--------------------------------------------------------------------------------------------------|
| **``tempest rm``**                   | remove the current directory from the target list                                                 |
| **``tempest rm 0``**                 | remove the target with index 0 from the target list                                               |
| **``tempest rm 1-3``**               | remove the targets with the index 1, 2 and 3                                                      |
| **``tempest rm *``**                 | remove all the targets                                                                            |
| **``tempest rm /tmp``**              | remove the target for the path: ``/tmp``                                                          |
| **``tempest rm 1-3 -o``**            | remove the targets for the matching directories to the indexes 1, 2 and 3                         |
| **``tempest rm 1-3 --origin``**       | same effect of using ``-o``                                                                       |


:warning: When using ``--origin`` or ``-o`` flag, the original files and directories will be removed from your system, so be careful.
:warning: After the use of ``rm``, Indexes numbers might change. Keep an eye on this. :eyes:
