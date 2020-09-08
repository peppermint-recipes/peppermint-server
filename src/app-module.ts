import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';

import { RecipeModule } from './recipe/recipe-module'
import { ConfigModule } from './config/config-module';
import { ConfigService } from './config/config-service';

@Module({
    imports: [RecipeModule,
        TypeOrmModule.forRootAsync({
            imports: [ConfigModule],
            useExisting: ConfigService,
        }),
    ],
    controllers: [],
    providers: [],
})
export class AppModule {}
