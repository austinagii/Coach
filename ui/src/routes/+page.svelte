<script>
   import Header from '../components/header/header.svelte'
  import Chat from '../components/chat/Chat.svelte'

</script>

<div class="page-wrapper">
  <Header />
  <main>
    <div class="chat-wrapper">
      <Chat />
    </div>
  </main>
</div>

<style>
  .page-wrapper {
    display: flex;
    flex-direction: column;
    height: 100vh;
  }

  main {
    height: 100%;
    padding: 0% 5%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
  }

  .chat-wrapper {
    width: 75%;
    height: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
  }
</style>
