We need to copy files to/from the filesystem that Godot uses to store it on our server or load them. To do that you have to understand how Emscripten virtual filesystem works and how to interact with it form your JavaScript.

Official Documentation:
https://emscripten.org/docs/porting/files/file_systems_overview.html

Log: 10.09.2024
To access WASM filesystem from JavaScript on the same page that hosts Godot WASM module you need to either export FS runtime module with `-sEXPORTED_RUNTIME_METHODS=['FS']` (then it will be availible on `Godot` module or `Engine` instances) or access the corresponding entires in IndexedDB, despite not showing up in the Chrome dev tools properly Godot mounts all of its user directories as persistent IndexedDB volumes and get calls to them work. Directory keys have not `contents` while file keys return `contents` in a int8Array as proper. I am still not sure how to list the directory contentents with IndexedDB.

We can get all of the IndexedDB contents with this code:
```js
const request = indexedDB.open("/home/web_user");
let files = [];

request.onsuccess = (event) => {
    let db = event.target.result
    let t = db.transaction(["FILE_DATA"])
    let objs = t.objectStore("FILE_DATA")
    objs.openCursor().onsuccess = (event) => {
      const cursor = event.target.result;
      if (cursor) {
        files.push(cursor.key);
        cursor.continue();
      } else {
        console.log(`Got all files: ${files}`);
      }
    }
}
```

Get all key-value pairs:
```js
const request = indexedDB.open("/home/web_user");

request.onsuccess = (event) => {
    let db = event.target.result
    let t = db.transaction(["FILE_DATA"])
    let objs = t.objectStore("FILE_DATA")
    objs.openCursor().onsuccess = (event) => {
            const cursor = event.target.result;
            if (cursor) {
                console.log(`Name for SSN ${cursor.key} is ${cursor.value.contents}`);
                cursor.continue();
          } else {
                console.log("No more entries!");
          }
    }
}
```

If `cursor.value.contents` is empty, then that entry is a directory.