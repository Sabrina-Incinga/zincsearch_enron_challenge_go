
<template> 
  <!--Modal-->
  <div :id="`modal-${mail.ID}`" class="modal opacity-0 pointer-events-none fixed w-full h-full top-0 left-0 flex items-center justify-center">
    <div :id="`modal-overlay-${mail.ID}`" class="modal-overlay absolute w-full h-full bg-gray-900 opacity-50"></div>
    
    <div class="modal-container bg-white w-11/12 mx-auto rounded shadow-lg z-50 overflow-y-auto">
      
      <div :class="`modal-close-${mail.ID} absolute top-0 right-0 cursor-pointer flex flex-col items-center mt-4 mr-4 text-white text-sm z-50`">
        <svg class="fill-current text-white" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 18 18">
          <path d="M14.53 4.53l-1.06-1.06L9 7.94 4.53 3.47 3.47 4.53 7.94 9l-4.47 4.47 1.06 1.06L9 10.06l4.47 4.47 1.06-1.06L10.06 9z"></path>
        </svg>
      </div>

      <!-- Add margin if you want to see some of the overlay behind the modal-->
      <div class="modal-content py-4 text-left px-6">
        <!--Title-->
        <div class="flex justify-between items-center pb-3">
          <p class="text-2xl font-bold"><strong>Message-ID: </strong> {{mail['Message-ID']}}</p>
          <div :class="`modal-close-${mail.ID} cursor-pointer z-50`">
            <svg class="fill-current text-black" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 18 18">
              <path d="M14.53 4.53l-1.06-1.06L9 7.94 4.53 3.47 3.47 4.53 7.94 9l-4.47 4.47 1.06 1.06L9 10.06l4.47 4.47 1.06-1.06L10.06 9z"></path>
            </svg>
          </div>
        </div>

        <!--Body-->
        <p><strong>From: </strong>{{ mail.From }}</p>
        <p><strong>To: </strong>{{ mail.To }}</p>
        <p><strong>Subject: </strong>{{ mail.Subject }}</p>
        <p><strong>Cc: </strong>{{ mail.Cc }}</p>
        <p><strong>Bcc: </strong>{{ mail.Bcc }}</p>
        <p><strong>Date: </strong>{{ mail.Date }}</p>
        <p><strong>Body: </strong>{{ mail.Body }}</p>


        <!--Footer-->
        <div class="flex justify-end pt-2">
          <button :class="`modal-close-${mail.ID} modal-close px-4 bg-indigo-500 p-3 rounded-lg text-white hover:bg-indigo-400`">Close</button>
        </div>
        
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { type Mail } from '@/types/mailTypes';
import { onMounted } from 'vue';

// eslint-disable-next-line vue/no-setup-props-destructure
const { mail }  = defineProps<{
  mail: Mail
}>()


onMounted(() => {
    var openmodal = document.getElementById(`modal-open-${mail.ID}`);
    console.log(openmodal)
    openmodal?.addEventListener('click', function(event: Event){
    event.preventDefault()
    toggleModal()
    })

    const overlay = document.getElementById(`modal-overlay-${mail.ID}`)
    overlay?.addEventListener('click', toggleModal)

    var closemodal = document.querySelectorAll(`.modal-close-${mail.ID}`)
    for (var i = 0; i < closemodal.length; i++) {
      closemodal[i].addEventListener('click', toggleModal)
    }

    document.onkeydown = function( evt: KeyboardEvent) {
        var isEscape = false
        if ("key" in evt) {
            isEscape = (evt.key === "Escape" || evt.key === "Esc")
            console.log(isEscape)
        } else {
            isEscape = evt.keyCode === 27
        }
        if (isEscape && document.body.classList.contains('modal-active')) {
            toggleModal()
        }
    };
    function toggleModal () {
        const body = document.querySelector('body')
        const modal = document.getElementById(`modal-${mail.ID}`)
        modal?.classList.toggle('opacity-0')
        modal?.classList.toggle('pointer-events-none')
        body?.classList.toggle('modal-active')
    }
});  


</script>
