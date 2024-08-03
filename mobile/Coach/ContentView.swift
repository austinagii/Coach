import SwiftUI
import os.log

// API URLs
let CHAT_API_URL = "http://localhost:8080/chats"

// Logger
let log = os.Logger()


// Main view
struct MainView: View {
    @State private var assistantMessage = ""
    @State private var userMessage = ""
    @State private var name = ""
    @State private var isLoading = false
    @State private var userId = ""
    @State private var chatId = ""
    @State private var errorOccurred = false
    @State private var errorMessage = ""
    
    var body: some View {
        VStack(spacing: 20) {
            if isLoading {
                Text("Loading...")
            } else {
                if errorOccurred {
                    Text(errorMessage)
                        .foregroundColor(   .red)
                } else if userId.isEmpty {
                    LabeledInput(label: "Enter your name", placeholder: "Your Name", text: $name)
                    Button(action: submitAction) {
                        Text("Submit")
                    }
                } else {
                    ChatView(userId: userId, chatId: chatId, initialMessage: assistantMessage)
                }
            }
        }
        .padding(20)
    }
    
    /// Handles the submission action by creating a user and a chat.
    private func submitAction() {
        Task {
            await handleSubmit()
        }
    }
    
    /// Handles the user and chat creation process.
    private func handleSubmit() async {
        if !name.isEmpty {
            isLoading = true
            defer { isLoading = false }
            
            guard let user = await createUser(name) else {
                setError("Failed to create user. Please try again.")
                return
            }
            userId = user.id
            
            guard let chat = await createChat(userId) else {
                setError("Failed to create chat. Please try again.")
                return
            }
            chatId = chat.id
            assistantMessage = chat.content
        }
    }
    
    /// Sets an error message and indicates an error occurred.
    private func setError(_ message: String) {
        errorOccurred = true
        errorMessage = message
    }
}

#Preview {
    MainView()
}
