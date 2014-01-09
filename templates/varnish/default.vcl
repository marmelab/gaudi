[[range (.Container.GetCustomValue "backends")]]
backend [[.]] {
    .host = ${[[ . | ToUpper ]]_PORT_[[ ($.Maestro.GetContainer . ).GetFirstPort ]]_TCP_ADDR};
    .port = ${[[ . | ToUpper ]]_PORT_[[ ($.Maestro.GetContainer . ).GetFirstPort ]]_TCP_PORT};
}
[[end]]
