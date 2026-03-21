# AppImageTool.go

A stripped down version of [the original appimagetool](https://github.com/AppImage/appimagetool) written in go.

This is going to be used in the Desktop.js CLI to allow for cross platform bundling of linux app versions.

- [x] Transformation of AppDir directory to squashfs filesystem
- [x] Download of appimage engine according to specified platform
- [x] Embed md5 integrity check
- [ ] Implement .upd_info to allow for incremental updates using zsync
- [ ] WIP: Implement signing using pgp keys
