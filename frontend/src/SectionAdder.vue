<script setup>
import { ref } from 'vue';
import { sections } from './globals';

const user_input = ref("")

function add_section() {
    console.log(user_input.value)
    console.log(sections.value)
    const tmp = user_input.value.split(' ')
    sections.value.push({
        name: tmp[0],
        link: tmp[1]
    })
}
</script>

<template>
    <form @submit.prevent="add_section">
        <input v-model="user_input" />
        <button type="submit">Add Section</button>
    </form>
</template>

<style lang="css" scoped>
form {
    display: block;
    margin: 1em 0;
}
</style>