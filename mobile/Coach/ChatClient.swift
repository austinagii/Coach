import Foundation

struct Message: Codable {
    
}
///
struct Chat: Codable {
    let id: String
    let messages: [String]
}

struct NewChatRequest: Codable {
    let userId: String
    let task: ChatTask
    
    enum CodingKeys: String, CodingKey {
        case userId = "user_id"
        case task
    }
}

struct ChatTask: Codable {
    let objective: String
}


/// Creates a new chat for the given user ID.
///
/// - Parameter userId: The ID of the user for whom to create a chat
/// - Returns: An optional Chat representing the newly created chat
func createChat(_ userId: String) async -> Chat? {
    let task = ChatTask(objective: "goal_creation")
    let newChatRequest = NewChatRequest(userId: userId, task: task)
    guard let url = URL(string: CHAT_API_URL) else {
        log.error("Could not create URL for chat API")
        return nil
    }
    
    guard let chat: Chat = await performHttpPostRequest(url: url, requestBody: newChatRequest) else {
        log.warning("Failed to create new chat")
        return nil
    }
    
    log.info("Chat created successfully")
    return chat
}

func sendChatMessage() async {
    let userMessage = UserMessage(userId: userId, text: userMessage)
    guard let url = URL(string:"\(CHAT_API_URL)/\(chatId)/messages") else {
        return
    }
    
    guard let assistantMessage: AssistantMessage = await performHttpPostRequest(url: url, requestBody: userMessage) else {
        return
    }
    
    self.userMessage = ""
    self.assistantMessage = assistantMessage.content
}
