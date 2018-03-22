
<a href="https://asciinema.org/a/cVCKdVX9lOeJEgBpPbbFrJX5f" target="_blank"><img src="https://asciinema.org/a/cVCKdVX9lOeJEgBpPbbFrJX5f.png" /></a>

### [](#head-config-1)Set up configuration:
:wrench: Configuration with **TEMPest** is done with the ``set`` command.  
```bash
tempest set
```
  

### [](#head-config-2)Access configuration:
:wrench: Access to configuration is done with the ``get`` command.  
```bash
tempest get
```

  
**Possible flags:**

| Long Flag       | Short Flag  | Description                                                                   |
|:----------------|:------------|:------------------------------------------------------------------------------|
| \-\-age         | -a          | Set the age of the targets. age = 3 means that if a file haven't been modified for 3 days, it will be deleted. |
| \-\-daemon-mode | -d          | ``true`` if you want it to run as daemon, ``false`` otherwise.                |
  
:warning: *Subject to changes*  