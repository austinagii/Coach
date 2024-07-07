<script>
  import { onMount } from 'svelte';

  import LoadingIndicator from '../loading-indicator/loading-indicator.svelte';

  const userName = "Kadeem";
  const createUserUrl = "http://localhost:8080/users"
  const createChatUrl = "http://localhost:8080/chats"
  const sendChatMessageUrl = "http://localhost:8080/chats/:chat_id/messages"

  let userResponseForm;

  let userId = "";
  let chatId = "";
  let assistantText = "";
  let userText = "";
  let isAssistantResponseLoading = true;
  let isAssistantTextAnimating = true;

  async function initChatSession() {
    isAssistantResponseLoading = true;
    userId = await createUser(userName);
    let chat = await createChat(userId);
    chatId = chat.id;
    isAssistantResponseLoading = false;
    animateAssistantText(chat.content, 2000)
  }

  /**
   * Creates a new user.
   * @param {string} name - The new user's name.
   * @return {Promise<string>} The user's ID.
   */
  async function createUser(name) {
    const response = await fetch(createUserUrl, {
      method: "POST", 
      body: JSON.stringify({ name: userName })
    })
    if (!response.ok) {
      console.log(response.status);
    }
   
    let json = await response.json();
    return json.id;
  }

  /**
   * Creates a new chat with an assistant.
   * @param {string} userId - The ID of the user starting the chat.
   * @return {Promise<any>} The ID of the chat.
   */
  async function createChat(userId) {
    const response = await fetch(createChatUrl, {
      method: "POST", 
      body: JSON.stringify({ 
        user_id: userId, 
        task: {
          objective: "goal_creation"
        } 
      })
    })
    if (!response.ok) {
      console.log(response.status);
    }
   
    let json = await response.json();
    return json;
  }

  async function animateAssistantText(text, duration) {
    if (text == null || text.length == 0) {
      return
    }
    // userResponseForm.classList.remove("fading-in");

    isAssistantTextAnimating = true;
    assistantText = "";
    let interval = Math.round(duration / text.length)
    let index = 0;
    function animate() {
      if (index < text.length) {
        assistantText += text.charAt(index);
        index++;
        setTimeout(animate, interval)
      } else {
        isAssistantTextAnimating = false;
        // userResponseForm.classList.add("fading-in");
      }
    }
    animate();
  }

  onMount(initChatSession);

  async function sendChatResponse(event) {
    event.preventDefault();

    assistantText = "";
    let url = sendChatMessageUrl.replace(":chat_id", chatId);
    isAssistantResponseLoading = true;
    let response = await fetch(url, {
        method: "POST",
        body: JSON.stringify({ 
            user_id: userId,
            text: userText
        })
    })
    let json = await response.json();

    isAssistantResponseLoading = false;
    userText = ""; 
    animateAssistantText(json.content, 2000)
  }
</script>

<section class="chat">
    {#if isAssistantResponseLoading}
        <LoadingIndicator />
    {:else}
        <p class="assistant-text">{assistantText}</p>
        <form on:submit={sendChatResponse} bind:this={userResponseForm} class="hidden {!isAssistantTextAnimating ? 'fading-in' : ''}">
            <input class="user-text" type="text" placeholder="Type your message here" bind:value={userText}/>
            <button type="submit">-></button>
        </form>
    {/if}
</section>

<style>
    :root {
        --font-size: 1.5em;
    }

    @keyframes blink {
        0% {opacity: 0%}
        50% {opacity: 100%}
        100% {opacity: 0%}
    }

    @keyframes fade-in {
        0% {opacity: 0%}
        100% {opacity: 100%}
    }

    .hidden {
        opacity: 0%;
    }

    .fading-in {
        animation: fade-in 3s forwards;
    }

    .chat {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 1em;
        height: 100%;
        width: 100%;
        font-family: 'Montserrat', sans-serif;
    }

    .assistant-text {
        font-size: var(--font-size);
    }

    .user-text {
        font-size: var(--font-size);
        margin-top: 1em;
        background-color: #000; 
        border: none;
        border-bottom: 1px solid white;
        color: white; 
    }

    .user-text::placeholder {
        color: #aaa;
    }

    .user-text:focus {
        outline: none;
    }
    
    form {
        opacity: 100%;
    }

    form button {
        display: none;
        font-size: var(--font-size);
        background-color: #000; 
        border: none;
        color: white; 
    }

    /* .user-text::after {
        content: "";
        position: absolute;
        width: calc(var(--font-size) * 0.3);
        height: calc(var(--font-size) * 0.8);
        background: white;
        margin-left: 2px;
        animation-name: blink;
        animation-iteration-count: infinite;
        animation-duration: 1s;
        animation-timing-function: ease-in-out;
    } */
</style>