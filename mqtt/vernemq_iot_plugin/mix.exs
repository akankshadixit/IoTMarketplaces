defmodule VernemqIoTPlugin.MixProject do
  use Mix.Project

  def project do
    [
      app: :vernemq_iot_plugin,
      version: "0.1.0",
      build_path: "./_build",
      config_path: "./config/config.exs",
      deps_path: "./deps",
      lockfile: "./mix.lock",
      elixir: "~> 1.9",
      start_permanent: Mix.env() == :prod,
      deps: deps(),
      releases: [
        vernemq_iot_plugin: [
          applications: [
            vernemq_iot_plugin: :permanent
          ],
          include_erts: false
        ]
      ]
    ]
  end

  # Run "mix help compile.app" to learn about applications.
  def application do
    [
      extra_applications: [:logger],
      mod: {VernemqIoTPlugin.Application, []},
      env: [
        vmq_plugin_hooks(),
        server_url: System.get_env("SERVER_URL", "http://172.17.0.1:8080"),
      ]
    ]
  end

  defp vmq_plugin_hooks do
    hooks = [
      {VernemqIoTPlugin, :auth_on_register, 5, []},
      {VernemqIoTPlugin, :on_register, 3, []},
      {VernemqIoTPlugin, :on_client_wakeup, 1, []},
      {VernemqIoTPlugin, :on_client_offline, 1, []},
      {VernemqIoTPlugin, :on_client_gone, 1, []},
      {VernemqIoTPlugin, :auth_on_subscribe, 3, []},
      {VernemqIoTPlugin, :on_subscribe, 3, []},
      {VernemqIoTPlugin, :on_unsubscribe, 3, []},
      {VernemqIoTPlugin, :auth_on_publish, 6, []},
      {VernemqIoTPlugin, :on_publish, 6, []},
      {VernemqIoTPlugin, :on_deliver, 4, []},
      {VernemqIoTPlugin, :on_offline_message, 5, []}
    ]

    {:vmq_plugin_hooks, hooks}
  end

  # Run "mix help deps" to learn about dependencies.
  defp deps do
    [
      # {:dep_from_hexpm, "~> 0.3.0"},
      # {:dep_from_git, git: "https://github.com/elixir-lang/my_dep.git", tag: "0.1.0"},
      # {:sibling_app_in_umbrella, in_umbrella: true},
      {:httpoison, "~> 1.8"},
      {:jason, "~> 1.2"}
    ]
  end
end
