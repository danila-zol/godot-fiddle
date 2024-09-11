Running this service is going to be quite expensive due to storage and bandwidth cost for retreiving objects.

We can use the following monetisation strategies:
- Premium tier:
	- Uploads up to 1GiB
	- C# support
	- Export a ready game for all platform instead of just downloading the project that you have to export yourself locally
	- Tutorials, possible integration with GDQuest
	- Features to monitor students, give assignments, and control a course for use in educational organisations. Godot Fiddle can become the preffered way to teach Godot at schools
- A research grant. This platform can be very interesting for educational organizations teaching kids how to code
- Run a kickstarter
- Ask for dontations
- Run fundrisers every half a year to sustain the site

However under any monetization strateges we are commited to keep the site open-source. Even with a premium tier option we will simply host a separate payment server (probably not open-source itself) and add a config option in the open-source backend `--paid-is-free` that makes any account have full premium features so that anyone can have the full premium features avaliable if they self-host our service. This can also be an attractive option for educational institutions to host on-premise. Hopefully they at least contribute code back or donate to the original project. Under no circumstances do I want this service to become non-OSS as it is build from and for the open-source community of Godot.