#include <lsp.hxx>
#include <server.hxx>

void runLSP()
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
        if (method == "textDocument/completion")
        {
            handleCompletion(request);
        }
        else if (method == "initialize")
        {
            json response;
            response["jsonrpc"] = "2.0";
            response["id"] = request["id"];
            response["result"] = {
                {"capabilities", {{"completionProvider", {{"resolveProvider", false}}}}}};
            sendMessage(response);
        }
    }
}