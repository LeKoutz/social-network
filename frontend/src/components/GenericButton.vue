<template>
    <button>
        <div
            class="flex flex-row align-items-center gap-6 background-color-input padding-block-start-6px padding-block-end-6px padding-inline-start-10px padding-inline-end-10px"
        >
            <slot name="icon"></slot>

            <span class="font-weight-500 text-color-accent font-size-13px"><slot></slot></span>
        </div>
    </button>
</template>
