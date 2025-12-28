#include <iostream>
#include <token.hxx>
#include <lsp.hxx>

void LSPServer::runLSP()
{
    while (true)
    {
        std::string msg = readMessage();
        if (msg.empty())
            continue;

        json request = json::parse(msg, nullptr, false);
        if (request.is_discarded())
            continue;

        std::string method = request.value("method", "");
        if (method == "initialize")
            handleInitialize(request);
        else if (method == "textDocument/completion")
            handleCompletion(request);
        else if (method == "shutdown")
            handleShutDown(request);
        else if (method == "exit")
            break;
    }
}

int LSPServer::readContentLength()
{
    std::string line;
    int contentLength = 0;
    while (std::getline(std::cin, line) && line != "\r" && !line.empty())
    {
        if (line.substr(0, 15) == "Content-Length:")
            contentLength = std::stoi(line.substr(15));
    }
    return contentLength;
}

std::string LSPServer::readMessage()
{
    int len = readContentLength();
    std::string msg(len, '\0');
    std::cin.read(&msg[0], len);
    return msg;
}

void LSPServer::sendMessage(const json &response)
{
    std::string s = response.dump();
    std::cout << "Content-Length: " << s.size() << "\r\n\r\n"
              << s;
    std::cout.flush();
}

void LSPServer::handleCompletion(const json &request)
{
    json response;
    response["jsonrpc"] = "2.0";
    response["id"] = request["id"];
    response["result"] = json::array();

    for (const auto &kv : keywords)
    {
        response["result"].push_back({{"label", kv.first},
                                      {"kind", 14},
                                      {"detail", kv.second}});
    }

    sendMessage(response);
}

void LSPServer::handleInitialize(const json &request)
{
    json response;
    response["jsonrpc"] = "2.0";
    response["id"] = request["id"];
    response["result"] = {
        {"capabilities", {{"completionProvider", {{"resolveProvider", false}}}}}};
    sendMessage(response);
}

void LSPServer::handleShutDown(const json &request)
{
    json response;
    response["jsonrpc"] = "2.0";
    response["id"] = request["id"];
    response["result"] = nullptr;
    sendMessage(response);
}