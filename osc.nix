{
  config,
  lib,
  pkgs,
  ...
}:
with lib; let
  cfg = config.services.osc;
in {
  #implementation
  options = {
    services.osc = {
      enable = mkEnableOption (lib.mdDoc "onesinglecounter application");
      ip = mkOption {
        default = "127.0.0.1";
        type = with types; uniq str;
      };
      port = mkOption {
        default = 5588;
        type = with types; uniq int;
      };
    };
  };
  # systemd service
  config = mkIf cfg.enable {
    # environment.etc."osc/firehol.conf".text = confFile;
    systemd.services.osc = {
      description = "onesinglecounter application";
      after = ["multi-user.target"];
      wants = ["network-online.target"];
      before = ["shutdown.target"];
      conflicts = ["shutdown.target"];
      restartIfChanged = true;
      serviceConfig = {
        Type = "oneshot";
        RemainAfterExit = true;
        ExecStart = "${pkgs.osc}/bin/osc --ip ${cfg.ip} --port ${cfg.port}";
        # ExecStop = "${pkgs.osc}/bin/fireqos stop";
        # ExecReload = "${pkgs.osc}/bin/firehol start";
      };
    };
  };
}
