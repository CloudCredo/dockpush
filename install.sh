go install

cf uninstall-plugin DockPush

cf install-plugin $GOPATH/bin/dockpush
