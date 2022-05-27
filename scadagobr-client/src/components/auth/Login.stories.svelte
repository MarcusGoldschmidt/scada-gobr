<script>
    import {Meta, Story, Template} from '@storybook/addon-svelte-csf';
    import Login from "./Login.svelte";
</script>

<Meta
        title="Button"
        component={Login}
        argTypes={{
    label: { control: 'text' },
    primary: { control: 'boolean' },
  }}
/>


<Template let:args>
    <Login {...args} />
</Template>

<!-- 👇 Each story then reuses that template -->

<Story
    name='Primary'
    args={{
        primary: true,
        label: 'Button'
    }}
/>