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

<!-- https://asciinema.org/a/bSMNAQieFKVbU4xXYARen1zbJ -->