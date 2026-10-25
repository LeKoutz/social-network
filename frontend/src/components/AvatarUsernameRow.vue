<script setup>
const props = defineProps({
    username: { type: String, default: null },
});
</script>

<template>
    <div
        class="display-flex position-relative width-28px height-28px background-7fa32eFF border-start-start-radius-14px border-end-end-radius-14px border-end-start-radius-14px border-start-end-radius-14px overflow-hidden align-items-center align-content-stretch justify-content-center flex-direction-row flex-wrap-nowrap flex-shrink-0"
    >
        <div
            class="display-inline-flex flex-direction-column justify-content-inherit margin-right-1px vertical-align-top"
        >
            <p class="font-size-0 line-height-1-2 margin-0 text-align-left" dir="ltr">
                <span
                    class="color-rgba-31-34-32-1 text-transform-none line-break-auto overflow-wrap-initial white-space-pre font-size-12px text-rendering-geometricPrecision caret-color-rgba-31-34-32-1 text-decoration-none letter-spacing-0px font-family-Inter font-style-normal font-weight-700"
                >
                    {{ username[0].toUpperCase() }}
                </span>
            </p>
        </div>
    </div>
    <div class="height-17px flex-shrink-0">
        <div>
            <div class="display-flex white-space-break-spaces align-items-flex-start">
                <div
                    class="display-inline-flex flex-direction-column justify-content-inherit margin-right-1px vertical-align-top"
                >
                    <p class="font-size-0 line-height-1-2 margin-0 text-align-left" dir="ltr">
                        <span
                            class="color-rgba-182-227-74-1 text-transform-none line-break-auto overflow-wrap-initial white-space-pre font-size-14px text-rendering-geometricPrecision caret-color-rgba-182-227-74-1 text-decoration-none letter-spacing-0px font-family-Inter font-style-normal font-weight-500"
                        >
                            {{ username }}
                        </span>
                    </p>
                </div>
            </div>
        </div>
    </div>
</template>
