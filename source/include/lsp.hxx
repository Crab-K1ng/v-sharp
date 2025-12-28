#pragma once

#include <iostream>
#include <token.hxx>
#include <nlohmann/json.hpp>

using json = nlohmann::json;

struct CompletionItem
{
    std::string label;
    std::string detail;
};

int readContentLength();
std::string readMessage();
void sendMessage(const json &response);
void handleCompletion(const json &request);