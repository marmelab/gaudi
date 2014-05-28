[[ $probeUrl := (.Container.GetCustomValue "probe_url" "/") ]]

[[range (.Container.GetCustomValue "backends")]]
[[ $port := ($.Collection.Get . ).GetFirstLocalPort ]]
backend [[.]] {
    .host = "${[[ . | ToUpper ]]_PORT_[[ $port ]]_TCP_ADDR}";
    .port = "${[[ . | ToUpper ]]_PORT_[[ $port ]]_TCP_PORT}";
    .probe = {
        .url = "[[ $probeUrl ]]";
        .interval = 5s;
        .timeout = 1 s;
        .window = 5;
        .threshold = 3;
      }
}
[[end]]

director loadBalancer round-robin {
[[range (.Container.GetCustomValue "backends")]]
        {
                .backend = [[.]];
        }
[[end]]
}

sub vcl_recv {
	set req.backend = loadBalancer;
}
