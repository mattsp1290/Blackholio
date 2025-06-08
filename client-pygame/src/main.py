"""
Blackholio Pygame Client
Entry point for the Python/Pygame client for Blackholio multiplayer circle game.

This client connects to a SpacetimeDB server and allows players to play alongside
Unity clients in the same game world.
"""

import sys
import logging
import argparse
import pygame

from game.game_manager import GameManager


def setup_logging(debug: bool = False) -> None:
    """
    Set up logging configuration.
    
    Args:
        debug: Whether to enable debug logging
    """
    log_level = logging.DEBUG if debug else logging.DEBUG  # Force DEBUG for now
    log_format = '%(asctime)s - %(name)s - %(levelname)s - %(message)s'
    
    logging.basicConfig(
        level=log_level,
        format=log_format,
        handlers=[
            logging.StreamHandler(sys.stdout),
            logging.FileHandler('blackholio_pygame.log')
        ]
    )
    
    # Suppress some verbose libraries
    logging.getLogger('urllib3').setLevel(logging.WARNING)
    logging.getLogger('websockets').setLevel(logging.WARNING)


def initialize_pygame() -> bool:
    """
    Initialize Pygame and set up the basic display.
    
    Returns:
        bool: True if initialization successful, False otherwise
    """
    try:
        pygame.init()
        
        # Try to initialize audio mixer, but don't fail if it doesn't work
        try:
            pygame.mixer.init()
            logging.info("Audio mixer initialized successfully")
        except pygame.error as audio_error:
            logging.warning(f"Failed to initialize audio mixer: {audio_error}")
            logging.info("Continuing without audio support")
        
        # Set up display - GameManager will take over from here
        screen = pygame.display.set_mode((1024, 768))
        pygame.display.set_caption("Blackholio - Pygame Client")
        
        logging.info("Pygame initialized successfully")
        return True
        
    except pygame.error as e:
        logging.error(f"Failed to initialize Pygame: {e}")
        return False


def parse_arguments() -> argparse.Namespace:
    """
    Parse command line arguments.
    
    Returns:
        Parsed arguments
    """
    parser = argparse.ArgumentParser(description="Blackholio Pygame Client")
    parser.add_argument(
        "--server", 
        default="ws://localhost:3000",
        help="SpacetimeDB server URL (default: ws://localhost:3000)"
    )
    parser.add_argument(
        "--debug",
        action="store_true",
        help="Enable debug logging"
    )
    parser.add_argument(
        "--player-name",
        default="Player",
        help="Default player name (default: Player)"
    )
    
    return parser.parse_args()


def main() -> int:
    """
    Main entry point for the Blackholio Pygame client.
    
    Returns:
        int: Exit code (0 for success, non-zero for error)
    """
    # Parse command line arguments
    args = parse_arguments()
    
    # Set up logging
    setup_logging(args.debug)
    
    logger = logging.getLogger(__name__)
    logger.info("Starting Blackholio Pygame Client...")
    logger.info(f"Server URL: {args.server}")
    logger.info(f"Default player name: {args.player_name}")
    
    # Initialize Pygame
    if not initialize_pygame():
        logger.error("Failed to initialize Pygame. Exiting.")
        return 1
    
    try:
        # Initialize and run game manager
        game_manager = GameManager(server_url=args.server)
        game_manager.run()
        
        logger.info("Game exited normally")
        return 0
        
    except KeyboardInterrupt:
        logger.info("Received keyboard interrupt, shutting down...")
        return 0
        
    except Exception as e:
        logger.error(f"Error during game execution: {e}", exc_info=True)
        return 1
        
    finally:
        pygame.quit()
        logger.info("Blackholio Pygame Client shut down.")


if __name__ == "__main__":
    sys.exit(main())
