# LuckyCat Frontend

A modern chat interface built with React, TypeScript, and Tailwind CSS.

## Project Structure

The project follows a modular architecture to ensure maintainability and readability:

```
src/
├── components/         # UI components
│   ├── ui/             # Shadcn UI components
│   ├── chat/           # Chat-specific components
│   │   ├── ChatHeader.tsx
│   │   ├── InputComponents.tsx
│   │   └── MessageComponents.tsx
│   ├── ChatInterface.tsx
│   └── LandingPage.tsx
├── hooks/              # Custom React hooks
│   └── useChat.ts      # Chat logic and state management
├── types/              # TypeScript type definitions
│   └── chat.ts         # Chat-related types
├── utils/              # Utility functions
│   └── chatUtils.ts    # Chat-specific utilities
├── lib/                # Library code
│   └── utils.ts        # General utilities
└── App.tsx             # Main application component
```

## Code Organization Principles

1. **Component Modularity**: Each component has a single responsibility and is kept small and focused.
2. **Type Safety**: TypeScript interfaces and types are used throughout the codebase.
3. **Custom Hooks**: Business logic is extracted into custom hooks for better separation of concerns.
4. **Utility Functions**: Common functionality is extracted into utility functions.
5. **Consistent Naming**: Clear and consistent naming conventions are used throughout the codebase.

## Key Components

### ChatInterface

The main component that orchestrates the chat experience. It uses the `useChat` hook for state management and renders the chat UI.

### MessageComponents

Contains components for rendering chat messages and message sections:

- `MessageItem`: Renders individual chat messages
- `MessageSectionComponent`: Organizes messages into sections

### InputComponents

Contains components for the chat input area:

- `ChatInputArea`: The main input container
- `InputButtons`: Action buttons for the chat interface

## Custom Hooks

### useChat

Manages all chat-related state and logic, including:

- Message state management
- Message streaming
- User input handling
- Chat actions (submit, toggle buttons)

## Utilities

### chatUtils

Contains utility functions for chat operations:

- `getAIResponse`: Generates AI responses
- `triggerVibration`: Handles device vibration
- `createUserMessage`: Creates user message objects
- `createSystemMessage`: Creates system message objects
- `adjustTextareaHeight`: Adjusts textarea height based on content

## Type Definitions

### chat.ts

Contains TypeScript interfaces and types for chat-related data:

- `Message`: Represents a chat message
- `MessageSection`: Represents a section of messages
- `ActiveButton`: Represents the active button state
- `MessageType`: Represents the type of message (user or system)
- `StreamingWord`: Represents a word being streamed in the UI

## Development

1. Clone the repository
2. Install dependencies: `npm install`
3. Start the development server: `npm run dev`

## Prerequisites

Before you begin, ensure you have the following installed:

- [Node.js](https://nodejs.org/) (v18 or higher recommended)
- [pnpm](https://pnpm.io/installation) (recommended package manager)

## Getting Started

Follow these steps to set up and run the project locally:

### 1. Clone the repository

```bash
git clone <repository-url>
cd luckycat-frontend
```

### 2. Install dependencies

```bash
pnpm install
```

### 3. Start the development server

```bash
pnpm dev
```

This will start the development server at [http://localhost:5173](http://localhost:5173) (or another port if 5173 is already in use).

## Available Scripts

In the project directory, you can run:

- `pnpm dev` - Starts the development server with hot-reload
- `pnpm build` - Builds the app for production to the `dist` folder
- `pnpm lint` - Runs ESLint to check for code quality issues
- `pnpm preview` - Locally preview the production build

## Technologies Used

- [React 19](https://react.dev/) - UI library
- [TypeScript](https://www.typescriptlang.org/) - Type safety
- [Vite](https://vitejs.dev/) - Build tool and development server
- [Tailwind CSS](https://tailwindcss.com/) - Utility-first CSS framework
- [shadcn/ui](https://ui.shadcn.com/) - Component library
- [Lucide React](https://lucide.dev/) - Icon library
- [ESLint](https://eslint.org/) - Code linting

## Customizing Components

This project uses shadcn/ui components which are fully customizable. You can modify the components in the `src/components/ui` directory to match your design requirements.

## Adding New Components

To add new shadcn/ui components, you can use the shadcn CLI:

```bash
npx shadcn-ui@latest add button
```

## Browser Support

This project supports modern browsers. For production environments, consider adding appropriate polyfills for older browser support if needed.

## License

[MIT](LICENSE)
