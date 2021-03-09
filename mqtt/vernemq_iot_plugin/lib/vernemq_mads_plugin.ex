defmodule VernemqIoTPlugin do

  def auth_on_register(_peer, {_mountpoint, clientid}, username, password, _clean_session?) do
    path = "/mqtt-auth-on-register"

    url = server_url() <> path
    params = %{clientid: clientid, username: username, password: password}
    params = Jason.encode!(params)
    headers = ["Accept": "Application/json; Charset=utf-8"]
    result = HTTPoison.post(url, params, headers)
    parse_auth_register(result)
    :ok
  end

  defp parse_auth_register({:ok, message}) do
    message = Jason.decode!(message.body)
    if message == "authenticated" do
      :ok
    else
      {:error, "authentication failed"}
    end
  end

  defp parse_auth_register({:error, message}) do
    IO.inspect(message)
    {:error, "some error occurred in registration"}
  end

  def on_register(_peer, {_mountpoint, clientid}, username) do
    IO.puts("*** on_register #{clientid} / #{username}")
    :ok
  end

  def on_client_wakeup({_mountpoint, clientid}) do
    IO.puts("*** on_client_wakeup #{clientid}")
    :ok
  end

  def on_client_offline({_mountpoint, clientid}) do
    IO.puts("*** on_client_offline #{clientid}")
    :ok
  end

  def on_client_gone({_mountpoint, clientid}) do
    IO.puts("*** on_client_gone #{clientid}")
    :ok
  end

  # Subscribe flow
  def auth_on_subscribe(_username, {_mountpoint, clientid}, topics) do
    IO.puts("*** auth_on_subscribe #{clientid}")
    {:ok, topics}
  end

  def on_subscribe(_username, {_mountpoint, clientid}, _topics) do
    IO.puts("*** on_subscribe #{clientid}")
    :ok
  end

  def on_unsubscribe(_username, {_mountpoint, clientid}, _topics) do
    IO.puts("*** on_unsubscribe #{clientid}")
    :ok
  end

  # Publish flow
  def auth_on_publish(_username, {_mountpoint, clientid}, _qos, topic, payload, _flag) do
    IO.puts("*** auth_on_publish #{clientid} / #{topic} / #{payload}")
    {:ok, payload}
  end

  def on_publish(_username, {_mountpoint, clientid}, _qos, topic, payload, _retain?) do
    IO.puts("*** on_publish #{clientid} / #{topic} / #{payload}")
    :ok
  end

  def on_deliver(_username, {_mountpoint, clientid}, topic, payload) do
    IO.puts("*** on_deliver #{clientid} / #{topic} / #{payload}")
    :ok
  end

  def on_offline_message({_mountpoint, clientid}, _qos, topic, payload, _retain?) do
    IO.puts("*** on_offline_message #{clientid} / #{topic} / #{payload}")
    :ok
  end

  defp server_url() do
    Application.get_env(:vernemq_iot_plugin, :server_url)
  end
end
