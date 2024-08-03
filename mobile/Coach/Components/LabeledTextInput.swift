import SwiftUI


/// Represents a text input with an associated text label
struct LabeledInput: View {
    let label: String
    let placeholder: String
    @Binding var text: String
    
    var body: some View {
        VStack(alignment: .leading) {
            Text(self.label)
            TextField(self.placeholder, text: $text)
                .padding(10)
                .overlay(
                    RoundedRectangle(cornerRadius: 5)
                        .stroke(Color.black.opacity(0.4), lineWidth: 2)
                )
        }
    }
}

#Preview {
    @State var firstName: String = ""
    
    return LabeledInput(label: "First Name", placeholder: "Enter your first name", text: $firstName)
}
