<script setup>
import { ref } from 'vue';
import { isBlank } from '../utils/validators.js';

// Reactive variable bound to the input field
const clientName = ref('');
const errorMessage = ref('');

// 1. Clear input (replaces javascript:void(0))
function clearForm() {
  clientName.value = '';
  errorMessage.value = '';
}

// 2. Validate blank input & open about:blank window
function generateReceipt() {
  if (isBlank(clientName.value)) {
    errorMessage.value = 'Please enter a valid client name.';
    return;
  }

  errorMessage.value = '';

  // 3. Open clean window and inject dynamic content
  const receiptWindow = window.open('about:blank', '_blank');
  if (receiptWindow) {
    receiptWindow.document.write(`<h1>Receipt for ${clientName.value}</h1>`);
    receiptWindow.document.close();
  }
}
</script>

<template>
  <!-- .prevent automatically calls e.preventDefault() -->
  <form @submit.prevent="generateReceipt">
    <input v-model="clientName" type="text" placeholder="Enter client name" />
    
    <p v-if="errorMessage" style="color: red;">{{ errorMessage }}</p>

    <button type="submit">Generate Receipt</button>

    <!-- @click.prevent stops the '#' link from navigating or scrolling -->
    <a href="#" @click.prevent="clearForm">Clear Form</a>
  </form>
</template>
