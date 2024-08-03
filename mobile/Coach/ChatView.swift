import SwiftUI

struct UserMessage: Codable {
    let userId: String
    let text: String
    
    enum CodingKeys: String, CodingKey {
        case userId = "user_id"
        case text
    }
}

struct AssistantMessage: Codable {
    let sender: String
    let content: String
}

struct ChatView: View {
    let userId: String
    let chatId: String
    let initialMessage: String
    
    @State private var assistantMessage = ""
    @State private var userMessage = ""
    
    var body: some View {
        VStack (alignment: .leading, spacing: 25) {
            Text(assistantMessage)
            HStack {
                TextField("Type your response here", text: $userMessage)
                Button(action: {
                    Task {
                        await sendChatMessage()
                    }
                }) {
                    Image(systemName: "paperplane")
                }
            }
        }
        .padding()
        .onAppear {
            print(initialMessage)
            assistantMessage = initialMessage
        }
    }
    
    
}

#Preview {
    ChatView(userId: "66a64cccad9ca17b0e4bde50", chatId: "66a64cddad9ca17b0e4bde51", initialMessage: "Hey there!")
}
