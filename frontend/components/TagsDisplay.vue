<script setup>
const tags = Array(50).fill("bug_games")
</script>

<template>
        <ul class=tags-panel>
            <li v-for="tag in tags"> {{tag}} </li>
        </ul>
</template>

<style lang="scss" scoped>
ul {
    height: 200px;
    display: flex;
    flex-wrap: wrap;
    flex-direction: column;
    list-style: none;

    li::before {
        content: "";
        width: 5px;
        height: 5px;
        border-radius: 50%;
        margin-right: 5px;
        display: inline-block;
        background-color: black;
        vertical-align: middle;
    }
}
</style>