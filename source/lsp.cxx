#include <lsp.hxx>

int readContentLength()
{
    std::string line;
    int contentLength = 0;
    while (std::getline(std::cin, line) && line != "\r" && !line.empty())
    {
        if (line.substr(0, 15) == "Content-Length:")
        {
            contentLength = std::stoi(line.substr(15));
        }
    }
    return contentLength;
}

std::string readMessage()
{
    int len = readContentLength();
    std::string msg(len, '\0');
    std::cin.read(&msg[0], len);
    return msg;
}

void sendMessage(const json &response)
{
    std::string s = response.dump();
    std::cout << "Content-Length: " << s.size() << "\r\n\r\n"
              << s;
    std::cout.flush();
}

void handleCompletion(const json &request)
{
    json response;
    response["jsonrpc"] = "2.0";
    response["id"] = request["id"];
    response["result"] = json::array();
    for (const auto &kv : keywords)
    {
        response["result"].push_back({{"label", kv.first},
                                      {"kind", 14},
                                      {"detail", "Keyword"}});
    }

    sendMessage(response);
}