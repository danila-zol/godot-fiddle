Needed patches:
- Remove the starting script to just start Godot. You just need to change the default html file for that.
- Disable useless warnings that are emitted in the web editor. Need to patch Godot source for that or better manipulate scons scripts used to build godot.
- Change the UI, remove most of the customization and advanced options. We are making a tool for demos, not complete fine-tuned games!

To compile the editor with our patches I will write a docker script. Thankfully Godot also provides web editor builds on GitHub releases, so you can alternatively just download that in a deploy script. It will not work if we apply any source patches though.
