defmodule VernemqIoTPlugin do

  def auth_on_register(_peer, {_mountpoint, clientid}, username, password, _clean_session?) do
    path = "/mqtt-auth-on-register"
    url = server_url() <> path
    params = %{actorid: username, token: password}
    params = Jason.encode!(params)
    headers = ["Accept": "*/*", "Content-Type": "application/json"]
    result = HTTPoison.post(url, params, headers)

    parse_auth_register(result)
  end

  defp parse_auth_register({:ok, message}) do
    message = Jason.decode!(message.body)
    if message["status"] == "success" do
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
    [{[topic| _], _qos}|_] = topics
    path = "/mqtt-auth-on-subscribe"
    url = server_url() <> path
    params = %{buyerid: clientid, streamid: topic}
    params = Jason.encode!(params)
    headers = ["Accept": "*/*", "Content-Type": "application/json"]
    result = HTTPoison.post(url, params, headers)

    parse_auth_on_subscribe(result)
  end

  defp parse_auth_on_subscribe({:ok, message}) do
    message = Jason.decode!(message.body)
    if message["status"] == "success" do
      :ok
    else
      {:error, "authentication failed"}
    end
  end

  defp parse_auth_on_subscribe({:error, message}) do
    IO.inspect(message)
    {:error, "some error occurred in registration"}
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
  def auth_on_publish(_username, {_mountpoint, clientid}, _qos, topics, payload, _flag) do
    [topic | _] = topics
    path = "/mqtt-auth-on-publish"
    url = server_url() <> path
    params = %{sellerid: clientid, streamid: topic}
    params = Jason.encode!(params)
    headers = ["Accept": "*/*", "Content-Type": "application/json"]
    result = HTTPoison.post(url, params, headers)

    parse_auth_on_publish(result, payload)
  end

  defp parse_auth_on_publish({:ok, message}, payload) do
    message = Jason.decode!(message.body)
    if message["status"] == "success" do
      {:ok, payload}
    else
      {:error, "authentication failed"}
    end
  end

  defp parse_auth_on_publish({:error, _message}, _payload) do
    {:error, "some error occurred in registration"}
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
