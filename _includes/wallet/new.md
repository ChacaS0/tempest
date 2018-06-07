This section shows of to create easily new ``targets``.

<a href="https://asciinema.org/a/bSMNAQieFKVbU4xXYARen1zbJ" target="_blank"><img src="https://asciinema.org/a/bSMNAQieFKVbU4xXYARen1zbJ.png" /></a>

### [](#head-add-1)The \*\*\*\*\* we doing:
* **``tempest rm *``:** clean the existing **targets** list.  
* **``tempest list``:** check the current **target** list.  
* **``tempest new --target --autoGen /tmp``:** create a new directory in ``/tmp`` called ``temp.est`` using the full name of flag and register it as a new **target** for **TEMPest**.  
* **``tempest new -t /tmp/test-tempest``:** create a new directory in ``/tmp/test-tempest`` using the short name for the flag ``--target`` and register it as a new **target** for **TEMPest**.  

### [](#head-new-2)List of flags

| Long Flag       | Short Flag  | Description                                                                                              |
|:----------------|:------------|:---------------------------------------------------------------------------------------------------------|
| \-\-autoGen     | -a          | [To be used with other flags] Combine this flag with ``--target <path>...`` in order to generate a default name for the target   |
| \-\-target      | -t          | [String]&lt;paths&gt;... Create the directory and registers the fresh created directory as a target in TEMPest your computer.    |


<!-- ### [](#head-new-3)Explanations

| Example command                      | Description                                                                                       |
|:-------------------------------------|:--------------------------------------------------------------------------------------------------|
| **``tempest rm``**                   | remove the current directory from the target list                                                 |
| **``tempest rm 0``**                 | remove the target with index 0 from the target list                                               |
| **``tempest rm 1-3``**               | remove the targets with the index 1, 2 and 3                                                      |
| **``tempest rm *``**                 | remove all the targets                                                                            |
| **``tempest rm /tmp``**              | remove the target for the path: ``/tmp``                                                          |
| **``tempest rm 1-3 -o``**            | remove the targets for the matching directories to the indexes 1, 2 and 3                         |
| **``tempest rm 1-3 --origin``**      | same effect of using ``-o``                                                                       | -->


tempest list
tempest new --target --autoGen /tmp
tempest list

tempest new -t /tmp/test-tempest
tempest list

<!-- https://asciinema.org/a/bSMNAQieFKVbU4xXYARen1zbJ -->