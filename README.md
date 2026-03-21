# AppImageTool.go

A stripped down version of [the original appimagetool](https://github.com/AppImage/appimagetool) written in go.

This is going to be used in the Desktop.js CLI to allow for cross platform bundling of linux app versions.

- [x] Create a good CLI experience
- [ ] Make it a go package which is easy to be imported and work with
- [x] Transformation of AppDir directory to squashfs filesystem
- [x] Download of appimage engine according to specified platform
- [x] Embed md5 integrity check
- [ ] Implement .upd_info to allow for incremental updates using zsync
- [ ] WIP: Implement signing using pgp keys

This implementation however should already be enough to create a valid [App Image Type 2 Format](https://github.com/AppImage/AppImageSpec/blob/master/draft.md).

```shell
Usage of appimagetool:
  -arch string
        System Architecture on which the AppImage should run. Valid values are: x86_64, aarch64, i686, armhf (default "x86_64")
  -passphrase string
        (Optional) Passphrase of encrypted PGP key file. Only use if encrypted.
  -runtime-file string
        (Optional) Path of AppImage runtime which BECOMES the AppImage
  -sign-key string
        (Optional) Path of PGP private key file to sign the AppImage
```
