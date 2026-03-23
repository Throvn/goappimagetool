# AppImageTool.go

A stripped down version of [the original appimagetool](https://github.com/AppImage/appimagetool) written in go.

This is going to be used in the Desktop.js CLI to allow for cross platform bundling of linux app versions.

- [x] Create a good CLI experience
- [ ] Make it a go package which is easy to be imported and work with
- [x] Transformation of AppDir directory to squashfs filesystem
- [x] Download of appimage engine according to specified platform
- [ ] Check downloaded appimage engine against the supplied (hardcoded) public key, to verify integrity
- [x] Embed md5 integrity check
- [ ] Implement .upd_info to allow for incremental updates using zsync
- [ ] WIP: Implement signing using pgp keys

This implementation however should already be enough to create a valid [App Image Type 2 Format](https://github.com/AppImage/AppImageSpec/blob/master/draft.md).

```
Usage of appimagetool [<folder>.AppDir, ...]:
  --arch string
        System Architecture on which the AppImage should run. Valid values are: x86_64, aarch64, i686, armhf (default "x86_64")
  --passphrase string
        (Optional) Passphrase of encrypted PGP key file. Only use if encrypted.
  --runtime-file string
        (Optional) Path of AppImage runtime which is copied into in the AppImage
  --sign-key string
        (Optional) Path of PGP private key file to sign the AppImage
```

As a convenience method, there also exists a command which lets you create a pgp key pair.
The command is `mkkey` and as an argument it takes the email which is encoded into the key pair.
As the certificate name, the currently logged in system user is taken (It really is just a convenience method).

```shell
appimagetool mkkey email@example.com [--passphrase]
```

you can then supply the created path of private.asc to the `--sign-key` param of the default command.
A passphrase can still be set by supplying the `--passphrase` param during command invocation.
