#pragma once

#include <nlohmann/json.hpp>

using json = nlohmann::json;

struct LSPServer
{
    void runLSP();

private:
    int readContentLength();
    std::string readMessage();
    void sendMessage(const json &response);
    void handleInitialize(const json &request);
    void handleShutDown(const json &request);
    void handleCompletion(const json &request);
};