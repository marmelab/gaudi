[[range (.Container.GetCustomValue "backends")]]
[[ $port := ($.Collection.Get . ).GetFirstPort ]]
backend [[.]] {
    .host = "${[[ . | ToUpper ]]_PORT_[[ $port ]]_TCP_ADDR}";
    .port = "${[[ . | ToUpper ]]_PORT_[[ $port ]]_TCP_PORT}";
    .probe = {
        .url = "/";
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
