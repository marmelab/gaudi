[[range (.Container.GetCustomValue "backends")]]
backend [[.]] {
    .host = "${[[ . | ToUpper ]]_PORT_[[ ($.Collection.Get . ).GetFirstPort ]]_TCP_ADDR}";
    .port = "${[[ . | ToUpper ]]_PORT_[[ ($.Collection.Get . ).GetFirstPort ]]_TCP_PORT}";
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
