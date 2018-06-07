
In this section, we'll see how to add a ``target`` to **TEMPest**, shall we begin? :smile_cat:

<!-- <script src="https://asciinema.org/a/171117.js" id="asciicast-171117" class="asciiframe" async></script> -->
<a href="https://asciinema.org/a/171117" target="_blank"><img src="https://asciinema.org/a/171117.png" /></a>

### [](#head-add-1)The \*\*\*\*\* we doing:
- ``tempest list``: Listing available targets (checking if there are any already set).
- ``tempest add /tmp``: Add ``/tmp`` as a target for **TEMPest**.
- ``tempest add ~/Documents/temp/ ~/Downloads/temp``: Add those two as targets, so we can add as many as we want in one command, ending by a ``/`` or not.
- ``tempest add``: Add the current directory as a target for **TEMPest**.

> **\#Note:** You cannot add multiple times the same target, or a target that doesnt exist.  
  
> :point_right: Use ``tempest add --auto`` to look for all the ``temp.est`` directories and register them as **targets**.

