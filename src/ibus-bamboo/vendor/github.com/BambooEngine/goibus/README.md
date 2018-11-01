goibus - golang implementation of libibus
==

goibus implements the libibus bindings in golang. goibus can be used to create IBus engines aka develop custom input methods.

IBus is an Intelligent Input Bus. It provides full featured and user friendly input method user interface. It also may help developers to develop input method easily.

This library is little bit different than other libibus bindings/wrappers. Instead of wrapping `libibus c library` or `GOBject-Introspection`, it implements whole functionality by communicating over DBus IPC. Because of that it is a independent 100% pure golang library without any native dependencies.

####NB:
libibus has various classes that are not absolutely required for creating engines. This library only implements engine related classes. Some uncommon class/methods are also skipped for now. You can always implement those and send PR ;)

This table shows the current status of implementation.

libibus | - | goibus
--- | --- | ---
[IBusAttrList](http://ibus.github.io/docs/ibus-1.5/IBusAttrList.html) | :white_check_mark: | Implemented In `text.go`
[IBusAttribute](http://ibus.github.io/docs/ibus-1.5/IBusAttribute.html) | :white_check_mark: | Implemented In `text.go`
[IBusBus](http://ibus.github.io/docs/ibus-1.5/IBusBus.html) | :white_check_mark: | Implemented In `bus.go`
[IBusComponent](http://ibus.github.io/docs/ibus-1.5/IBusComponent.html) | :white_check_mark: | Implemented In `component.go`
[IBusConfig](http://ibus.github.io/docs/ibus-1.5/IBusConfig.html) | :red_circle: | Ignored, not implemented
[IBusConfigService](http://ibus.github.io/docs/ibus-1.5/IBusConfigService.html) | :red_circle: | Ignored, not implemented
[IBusEngine](http://ibus.github.io/docs/ibus-1.5/IBusEngine.html) | :white_check_mark: | Implemented In `engine.go`
[IBusEngineDesc](http://ibus.github.io/docs/ibus-1.5/IBusEngineDesc.html) | :white_check_mark: | Implemented In `engineDesc.go`
[IBusFactory](http://ibus.github.io/docs/ibus-1.5/IBusFactory.html) | :white_check_mark: | Implemented In `factory.go`
[IBusHotkeyProfile](http://ibus.github.io/docs/ibus-1.5/IBusHotkeyProfile.html) | :red_circle: | Ignored, not implemented
[IBusInputContext](http://ibus.github.io/docs/ibus-1.5/IBusInputContext.html) | :large_blue_circle: | Ignored, relevant inherited signals implemented in `Engine`
[IBusKeymap](http://ibus.github.io/docs/ibus-1.5/IBusKeymap.html) | :large_blue_circle: | Ignored for now, will implement
[IBusLookupTable](http://ibus.github.io/docs/ibus-1.5/IBusLookupTable.html) | :white_check_mark: | Implemented In `lookupTable.go`
[IBusObject](http://ibus.github.io/docs/ibus-1.5/IBusObject.html) | :white_check_mark: | Ignored, Parent/Interface class, relevant inherited signals implemented in `Engine`
[IBusObservedPath](http://ibus.github.io/docs/ibus-1.5/IBusObservedPath.html) | :red_circle: | Ignored, not implemented
[IBusPanelService](http://ibus.github.io/docs/ibus-1.5/IBusPanelService.html) | :red_circle: | Ignored, not implemented
[IBusPropList](http://ibus.github.io/docs/ibus-1.5/IBusPropList.html) | :white_check_mark: | Implemented In `property.go`
[IBusProperty](http://ibus.github.io/docs/ibus-1.5/IBusProperty.html) | :white_check_mark: | Implemented In `property.go`
[IBusProxy](http://ibus.github.io/docs/ibus-1.5/IBusProxy.html) | :red_circle: | Ignored, not implemented
[IBusRegistry](http://ibus.github.io/docs/ibus-1.5/IBusRegistry.html) | :red_circle: | Ignored, not implemented
[IBusSerializable](http://ibus.github.io/docs/ibus-1.5/IBusSerializable.html) | :white_check_mark: | Not needed in golang, All implemented classes are Serializable
[IBusService](http://ibus.github.io/docs/ibus-1.5/IBusService.html) | :white_check_mark: | Ignored, not needed. Parent/Interface class
[IBusText](http://ibus.github.io/docs/ibus-1.5/IBusText.html) | :white_check_mark: | Implemented In `text.go`


Installation
==

```
go get github.com/godbus/dbus
go get github.com/sarim/goibus
```

check `_example` directory for a sample engine and ~~ TODO:detailed tutorial ~~. Run the sample engine by `_example -standalone` to see it in action.
![sample engine](https://cloud.githubusercontent.com/assets/1235888/7563038/569ef518-f7fb-11e4-91af-2c2150199fe7.png)

License
==
**goibus** - golang implementation of libibus by **Sarim Khan**

Licensed under Mozilla Public License 1.1 ("MPL"), an open source/free software license.
